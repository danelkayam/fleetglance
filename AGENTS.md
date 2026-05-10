# Fleetglance

Fleetglance is a lightweight visual telemetry system for a self-hosted homelab fleet. It exists to solve a simple operational problem: seeing the short health state of several rack machines at a glance, on a dedicated display, without deploying a heavy observability stack.

The project uses a small pull-based model: each server runs an agent that exposes telemetry over HTTP, and a console polls those agents, converts results into events, and renders a compact Bubble Tea TUI.

## Terminology

- **Fleet**
  - The full set of monitored machines.
  - Defined by the console YAML config.
  - Current config version is `version: 1`.

- **Ship**
  - One monitored server, VM, or host in the fleet.
  - In config, each ship has a stable name and an agent base URL.
  - Example: `donnager: { url: http://donnager:9800 }`.

- **Agent**
  - The lightweight HTTP service running on each ship.
  - Collects local telemetry from the host and exposes it through `/api/telemetry`.
  - Current binary entrypoint: `cmd/agent`.
  - Current implementation package: `internal/agent`.

- **Console**
  - The local TUI application that displays fleet state.
  - Loads the fleet config, starts the telemetry engine, receives events, and renders the UI.
  - Current binary entrypoint: `cmd/console`.
  - Current implementation package: `internal/console`.

- **Telemetry**
  - The structured health data returned by an agent.
  - Current protocol fields include agent version, timestamp, uptime, CPU, memory, storage, and containers.
  - Defined in `internal/protocol/telemetry.go`.

- **Event**
  - A console-side message produced by the engine for a specific ship.
  - Contains the ship name, telemetry payload, or error.
  - Defined in `internal/console/engine/event.go`.

- **Protocol**
  - Shared JSON response and telemetry structures used by agent and console.
  - Lives under `internal/protocol`.
  - JSON field names are snake_case and must remain stable unless a protocol change is explicitly requested.

## Architecture

### Agent flow

1. `cmd/agent/main.go` loads runtime params from `.env` and environment variables.
2. It creates `agent.NewAgent(...)`.
3. `Agent.Start()` creates the Gin router and starts an HTTP server.
4. The router exposes:
   - `GET /healthz`
   - `GET /api/telemetry`
5. The telemetry handler calls `providers.TelemetryProvider.GetTelemetry()`.
6. The provider collects host telemetry and returns `protocol.Telemetry`.
7. The handler wraps the payload in `protocol.Response[T]`.

Current agent defaults:

- `PORT=9800`
- `DEBUG=false`
- `LOG_FORMAT=console`

Telemetry collection rules:

- CPU, memory, storage, uptime, and container data are collected locally.
- Docker container telemetry uses the Docker socket and may return `nil` if unavailable.
- Collector failures should warn and return a nil subsection where possible, not fail the whole telemetry response unless explicitly intended.
- Percent values are rounded to one decimal place.

### Console flow

1. `cmd/console/main.go` loads params from `.env`, environment variables, and flags.
2. Config path can come from:
   - `FLEETGLANCE_CONFIG_PATH`
   - `-f path/to/fleetglance.yaml`
3. The `-f` flag overrides `FLEETGLANCE_CONFIG_PATH`.
4. `config.LoadFleet(path)` parses YAML and applies defaults.
5. `console.NewConsole(fleet)` creates a single-use console runtime.
6. `Console.Start()` validates the fleet config, starts the engine, starts Bubble Tea, and forwards engine events to the UI.
7. `Console.Stop()` stops the engine and quits the Bubble Tea program.

Current console config defaults:

- `pull_interval: 5s`
- `timeout: 2s`
- maximum ships: `8`

Current config shape:

```yaml
version: 1

pull_interval: 5s
timeout: 2s

ships:
  donnager:
    url: http://donnager:9800
  nostromo:
    url: http://nostromo:9800
```

Important config rule:

- Ship URLs are agent base URLs only.
- Do not put `/api/telemetry` in config URLs.
- The console telemetry client appends `/api/telemetry` internally.

### Engine flow

1. The console engine converts configured ships into sorted internal `Ship` values.
2. Each ship gets its own telemetry stream goroutine.
3. A stream performs one immediate telemetry fetch, then repeats every `pull_interval` using a ticker.
4. Each fetch emits an `engine.Event` containing:
   - `ShipName`
   - `Telemetry`
   - `Error`
5. The engine closes the event channel after shutdown.
6. The engine is single-use. Starting it twice is invalid.

### UI flow

1. `ui.NewModel(fleet)` builds the initial Bubble Tea model from configured ships.
2. `TelemetryEventMsg` updates per-ship state.
3. A one-second tick updates the clock.
4. The view renders:
   - top fleet summary
   - ship panes
   - per-ship status, uptime, CPU, RAM, disk, and containers
5. Ship names are sorted for deterministic rendering.

### Shared protocol

Successful agent response:

```json
{
  "data": {
    "agent_version": "dev-unknown",
    "timestamp": "2026-05-06T12:00:00Z",
    "uptime_seconds": 123,
    "cpu": { "usage_percent": 12.5 },
    "memory": {
      "used_bytes": 1024,
      "total_bytes": 2048,
      "usage_percent": 50.0
    },
    "storage": {
      "used_bytes": 4096,
      "total_bytes": 8192,
      "usage_percent": 50.0
    },
    "containers": {
      "running": 2,
      "total": 3
    }
  }
}
```

Error response:

```json
{
  "data": null,
  "error": {
    "message": "failed"
  }
}
```

## Project structure

- `cmd/agent`
  - Agent process entrypoint.
  - Loads env params, initializes logging, starts/stops the agent.

- `cmd/console`
  - Console process entrypoint.
  - Loads env/flag params, loads fleet YAML, starts/stops the console.

- `internal/agent`
  - Agent runtime, router, handlers, and telemetry providers.

- `internal/agent/providers`
  - Host telemetry collectors.
  - Keep collectors small and focused by telemetry domain.

- `internal/agent/routers`
  - Gin router construction and router options.

- `internal/agent/routers/handlers`
  - HTTP handlers and error mapping.

- `internal/console`
  - Console lifecycle and Bubble Tea program ownership.

- `internal/console/config`
  - Fleet YAML loading, defaults, and validation.

- `internal/console/engine`
  - Pull-based telemetry engine and HTTP telemetry client.

- `internal/console/ui`
  - Bubble Tea model, messages, formatting, theme, and view rendering.

- `internal/protocol`
  - Shared API protocol types.

- `internal/version`
  - Build-time version metadata injected through ldflags.

- `docker/Dockerfile.agent`
  - Agent container image build.

## Tools

Use the Makefile. Do not invent parallel local commands unless there is a specific reason.

Common commands:

```sh
make format
make lint
make test
make build
make build-agent
make build-console
make clean
make mod-download
make mod-tidy
make build-docker-agent
make push-docker-agent-multi
```

Required validation before finishing any code change:

```sh
make format lint test
```

Build metadata is injected through Makefile ldflags:

- `VERSION`
- `COMMIT`
- `BUILT_AT`

Development tool binaries are installed under `bin-tools/` by Makefile targets. Build outputs go under `bin/`.

## Guidelines

### General rules

- Keep this project boring, explicit, and operationally obvious.
- Write boring, well-structured Go: simple where possible, composed where useful, and never a 500-line pile of unrelated rendering logic.
- Do not introduce architectural churn unless the task explicitly asks for it.
- Do not rename project terms unless explicitly asked.
- Use the established terms: fleet, ship, agent, console, telemetry, event, protocol.
- Do not convert this into a web UI, dashboard server, database-backed system, Prometheus stack, or generic observability framework unless explicitly asked.
- This is a lightweight self-hosted visual telemetry project for a small homelab fleet.

### Code style

- Write idiomatic Go.
- Use composition and small interfaces where they clarify boundaries.
- Do not collapse unrelated behavior into large procedural files.
- Avoid 500+ line grab-bag modules.
- Keep package responsibilities narrow.
- Keep public APIs minimal.
- Avoid global state except for simple build metadata or unavoidable library configuration.
- Return wrapped errors with useful context.
- Prefer deterministic ordering for maps when output, UI order, or tests depend on it.
- Keep JSON protocol fields stable and snake_case.

### Testing rules

- Any new behavior must have tests.
- Any changed behavior must have updated tests.
- Test both success and failure paths.
- For every test that validates a good value, add the opposite case where practical:
  - bad input
  - missing input
  - invalid config
  - error response
  - nil data
  - edge value
- Prefer table-driven tests for validators, formatters, protocol shape, and client error cases.
- Prefer `httptest` and small local fakes for HTTP and provider boundaries.
- Do not add heavy mocking frameworks unless explicitly requested.
- Do not rewrite production code only to make tests easier unless the task explicitly allows refactoring.
- Protocol tests must protect the JSON shape, especially snake_case field names.

### Agent rules

- The agent owns local telemetry collection only.
- The agent should expose simple HTTP endpoints.
- Keep `/api/telemetry` stable.
- Keep `/healthz` simple.
- If one telemetry subsection fails, prefer logging a warning and returning `nil` for that section over failing the whole response.
- Do not add authentication, persistence, service discovery, remote control, or scheduling unless explicitly requested.
- Do not change default port `9800` unless explicitly requested.

### Console rules

- The console owns config loading, telemetry pulling, event forwarding, and TUI rendering.
- The console must remain pull-based.
- Config ship URLs are base URLs. Do not require `/api/telemetry` in config.
- Keep `FLEETGLANCE_CONFIG_PATH` support.
- Keep `-f` overriding `FLEETGLANCE_CONFIG_PATH`.
- Keep console and engine lifecycle single-use unless explicitly asked to redesign lifecycle semantics.
- `Console.Start()` starts the UI and blocks until exit or stop.
- `Console.Stop()` should be safe to call more than once.

### UI rules

- The UI is Bubble Tea based.
- Keep UI state updates event-driven.
- Avoid direct telemetry fetching inside UI code.
- Keep formatting helpers tested.
- Keep ship ordering deterministic.
- Do not hard-code specific example ship names into theme or logic.
- If assigning ship colors, use palette/index-based assignment rather than name-based assignment.
- Maintain usability for small dedicated displays where possible.

### Protocol rules

- Shared protocol types live in `internal/protocol`.
- Keep the response envelope generic: `protocol.Response[T]`.
- Success responses use `data`.
- Error responses use `data: null` and `error.message`.
- Do not add fields to the public telemetry protocol without tests.
- Do not rename JSON fields without explicit approval.

### Dependency rules

- Keep dependencies minimal and practical.
- Do not add a dependency for trivial code.
- Existing important dependencies:
  - Gin for agent HTTP routing.
  - Bubble Tea and Lip Gloss for console TUI.
  - gopsutil for host telemetry.
  - caarlos0/env and godotenv for env loading.
  - validator for parameter validation.
  - yaml.v3 for fleet YAML parsing.
  - zerolog for logging.
- Do not switch YAML libraries, routers, TUI framework, or logging libraries unless explicitly requested.

### Git rules

- Do not amend git history.
- Do not run `git rebase`.
- Do not run `git reset --hard` unless explicitly instructed.
- Do not force-push.
- Do not make history-changing git operations.
- Do not create commits unless explicitly requested.
- It is acceptable to inspect git state with read-only commands such as `git status` and `git diff`.

### Change-control rules

- Do not change existing architecture, terminology, package layout, config shape, endpoints, or protocol unless the task explicitly asks for it.
- Do not “improve” adjacent areas while working on a narrow request.
- If a requested change conflicts with current architecture, point it out before making broad changes.
- If requirements are ambiguous, ask and verify unless the user explicitly asked for a best-effort implementation.
- Keep changes scoped to the requested area.
- Avoid recursive scopes of recursive scopes. Cut the work back to v1 before it becomes a journey to Mordor.

## Definition of done

A change is done only when:

1. The requested behavior is implemented.
2. Relevant tests were added or updated.
3. `make format lint test` passes.
4. Existing architecture and terms are preserved unless explicitly changed by the task.
5. No unrelated files were modified.
6. No git history was rewritten.
