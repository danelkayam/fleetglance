# Installation

Fleetglance is currently pre-release and in rapid development.

This guide describes the basic installation flow.

## Agent installation via Docker

Run the agent on each monitored machine.

Docker Compose example:

```yaml
services:
  fleetglance-agent:
    image: danelkayam/fleetglance-agent:latest
    container_name: fleetglance-agent
    ports:
      - "9800:9800"
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
    restart: unless-stopped
```

Verify that the agent is reachable:

```sh
curl http://localhost:9800/api/telemetry
```

## Console installation

Install the `fleetglance-console` binary on the machine that will run the display.

Recommended path:

```sh
sudo install -m 0755 fleetglance-console /usr/local/bin/fleetglance-console
```

Verify:

```sh
fleetglance-console --version
```

## Create configuration

Create the config directory:

```sh
mkdir -p ~/.config/fleetglance
```

Create the config file:

```sh
vim ~/.config/fleetglance/fleetglance.yaml
```

Example:

```yaml
version: 1

pull_interval: 5s
timeout: 2s

ships:
  donnager:
    url: http://donnager.example.local:9800
  nostromo:
    url: http://nostromo.example.local:9800
  rocinante:
    url: http://rocinante.example.local:9800
```

## Run the console

```sh
fleetglance-console -f ~/.config/fleetglance/fleetglance.yaml
```

## Verify from the console machine

From the machine running the console, verify that each agent is reachable:

```sh
curl http://donnager.example.local:9800/api/telemetry
```

If this does not work, the console will not be able to collect telemetry from that ship.
