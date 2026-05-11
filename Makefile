APP_NAME := fleetglance

GOOS ?= linux
GOARCH ?= amd64
PLATFORMS ?= linux/amd64,linux/arm64

BIN_DIR := $(CURDIR)/bin
BIN_TOOLS_DIR := $(CURDIR)/bin-tools
DIST_DIR := $(CURDIR)/dist

BIN_AGENT := $(BIN_DIR)/fleetglance-agent
SRC_AGENT := $(CURDIR)/cmd/agent/main.go
BIN_CONSOLE := $(BIN_DIR)/fleetglance-console
SRC_CONSOLE := $(CURDIR)/cmd/console/main.go

DIST_AGENT_LINUX_AMD64 := $(DIST_DIR)/fleetglance-agent_linux_amd64
DIST_AGENT_LINUX_ARM64 := $(DIST_DIR)/fleetglance-agent_linux_arm64
DIST_CONSOLE_LINUX_AMD64 := $(DIST_DIR)/fleetglance-console_linux_amd64
DIST_CONSOLE_LINUX_ARM64 := $(DIST_DIR)/fleetglance-console_linux_arm64

COVERAGE_PROFILE := $(CURDIR)/coverage.out

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
BUILT_AT ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

LDFLAGS := \
	-X '$(APP_NAME)/internal/version.Version=$(VERSION)' \
	-X '$(APP_NAME)/internal/version.Commit=$(COMMIT)' \
	-X '$(APP_NAME)/internal/version.BuiltAt=$(BUILT_AT)'


DOCKER_DIR := $(CURDIR)/docker

DOCKER_REPO ?= danelkayam
AGENT_IMAGE := $(DOCKER_REPO)/$(APP_NAME)-agent

.PHONY: \
	help build build-release clean test lint format update  \
	build-deps clean-deps mod-download mod-tidy \
	tools-install tools-clean archive \
	build-agent build-console \
	build-agent-linux-amd64 build-agent-linux-arm64 \
	build-console-linux-amd64 build-console-linux-arm64 \
	package-release build-docker-agent push-docker-agent-multi

help:
	@echo "Available commands:"
	@echo "  make build                   - Build local/dev agent and console binaries"
	@echo "  make build-agent             - Build the local/dev agent binary"
	@echo "  make build-console           - Build the local/dev console binary"
	@echo "  make build-release           - Build and package linux/amd64 and linux/arm64 release binaries"
	@echo "  make clean                   - Clean build artifacts"
	@echo "  make test                    - Run tests with coverage"
	@echo "  make lint                    - Run linting"
	@echo "  make format                  - Run code formatting"
	@echo "  make update                  - Update go modules"
	@echo "  make mod-download            - Download go module dependencies"
	@echo "  make mod-tidy                - Tidy go module dependencies"
	@echo "  make tools-install           - Install development tools"
	@echo "  make tools-clean             - Clean development tools"
	@echo "  make archive                 - Create code archive"
	@echo "  make build-docker-agent      - Build agent docker image"
	@echo "  make push-docker-agent-multi - Push multi-arch agent docker image to registry"

default: build

build: build-agent build-console

build-release: \
	build-agent-linux-amd64 \
	build-agent-linux-arm64 \
	build-console-linux-amd64 \
	build-console-linux-arm64 \
	package-release

build-agent:
	@echo "Building $(APP_NAME) agent..."
	@mkdir -p $(BIN_DIR)
	@CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) \
		go build -ldflags "$(LDFLAGS)" -o $(BIN_AGENT) $(SRC_AGENT)

build-console:
	@echo "Building $(APP_NAME) console..."
	@mkdir -p $(BIN_DIR)
	@CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) \
		go build -ldflags "$(LDFLAGS)" -o $(BIN_CONSOLE) $(SRC_CONSOLE)

build-agent-linux-amd64:
	@echo "Building $(APP_NAME) agent for linux/amd64..."
	@mkdir -p $(DIST_DIR)
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
		go build -ldflags "$(LDFLAGS)" -o $(DIST_AGENT_LINUX_AMD64) $(SRC_AGENT)

build-agent-linux-arm64:
	@echo "Building $(APP_NAME) agent for linux/arm64..."
	@mkdir -p $(DIST_DIR)
	@CGO_ENABLED=0 GOOS=linux GOARCH=arm64 \
		go build -ldflags "$(LDFLAGS)" -o $(DIST_AGENT_LINUX_ARM64) $(SRC_AGENT)

build-console-linux-amd64:
	@echo "Building $(APP_NAME) console for linux/amd64..."
	@mkdir -p $(DIST_DIR)
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
		go build -ldflags "$(LDFLAGS)" -o $(DIST_CONSOLE_LINUX_AMD64) $(SRC_CONSOLE)

build-console-linux-arm64:
	@echo "Building $(APP_NAME) console for linux/arm64..."
	@mkdir -p $(DIST_DIR)
	@CGO_ENABLED=0 GOOS=linux GOARCH=arm64 \
		go build -ldflags "$(LDFLAGS)" -o $(DIST_CONSOLE_LINUX_ARM64) $(SRC_CONSOLE)

package-release:
	@echo "Packaging release binaries..."
	@tar -C $(DIST_DIR) -czf $(DIST_AGENT_LINUX_AMD64).tar.gz $(notdir $(DIST_AGENT_LINUX_AMD64))
	@tar -C $(DIST_DIR) -czf $(DIST_AGENT_LINUX_ARM64).tar.gz $(notdir $(DIST_AGENT_LINUX_ARM64))
	@tar -C $(DIST_DIR) -czf $(DIST_CONSOLE_LINUX_AMD64).tar.gz $(notdir $(DIST_CONSOLE_LINUX_AMD64))
	@tar -C $(DIST_DIR) -czf $(DIST_CONSOLE_LINUX_ARM64).tar.gz $(notdir $(DIST_CONSOLE_LINUX_ARM64))

clean:
	@echo "Cleaning $(APP_NAME) build artifacts..."
	@rm -rf $(BIN_DIR)/* $(DIST_DIR)/* $(COVERAGE_PROFILE)

test:
	@echo "Running tests..."
	go test -v ./... -cover -coverprofile=$(COVERAGE_PROFILE)

lint: $(BIN_TOOLS_DIR)/golangci-lint \
	  $(BIN_TOOLS_DIR)/modernize
	@echo "Running linting..."
	@go vet ./...
	@$(BIN_TOOLS_DIR)/golangci-lint run ./...
	@$(BIN_TOOLS_DIR)/modernize ./...

format: $(BIN_TOOLS_DIR)/goimports
	@echo "Running code formatting..."
	@$(BIN_TOOLS_DIR)/goimports -w cmd/ internal/
	@go fmt ./...

update:
	@echo "Updating go modules..."
	@go get -u ./...
	@go mod tidy

mod-download:
	@echo "Downloading go module dependencies..."
	@go mod download

mod-tidy:
	@echo "Tidying go module dependencies..."
	@go mod tidy

tools-install:
	@echo "Installing development tools..."

tools-clean:
	@echo "Cleaning development tools..."
	@rm -rf $(BIN_TOOLS_DIR)/*

archive:
	@echo "Creating code archive..."
	@git archive --format=tar.gz -o $(APP_NAME)-$(VERSION).tar.gz HEAD

build-docker-agent:
	@echo "Building agent docker image..."
	@docker build -f $(DOCKER_DIR)/Dockerfile.agent $(CURDIR) \
		--build-arg TARGETOS=$(GOOS) \
		--build-arg TARGETARCH=$(GOARCH) \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT=$(COMMIT) \
		--build-arg BUILT_AT=$(BUILT_AT) \
		-t $(AGENT_IMAGE):$(VERSION) \
		-t $(AGENT_IMAGE):latest

push-docker-agent-multi:
	@echo "Building and pushing multi-arch agent docker image..."
	@docker buildx build \
		-f $(DOCKER_DIR)/Dockerfile.agent \
		--platform $(PLATFORMS) \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT=$(COMMIT) \
		--build-arg BUILT_AT=$(BUILT_AT) \
		-t $(AGENT_IMAGE):$(VERSION) \
		-t $(AGENT_IMAGE):latest \
		--no-cache \
		--provenance=false \
		--push \
		$(CURDIR)

$(BIN_TOOLS_DIR):
	@mkdir -p $(BIN_TOOLS_DIR)

$(BIN_TOOLS_DIR)/golangci-lint: | $(BIN_TOOLS_DIR)
	@echo "Installing golangci-lint..."
	@env GOBIN=$(BIN_TOOLS_DIR) go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

$(BIN_TOOLS_DIR)/modernize: | $(BIN_TOOLS_DIR)
	@echo "Installing modernize..."
	@env GOBIN=$(BIN_TOOLS_DIR) go install golang.org/x/tools/gopls/internal/analysis/modernize/cmd/modernize@latest

$(BIN_TOOLS_DIR)/goimports: | $(BIN_TOOLS_DIR)
	@echo "Installing goimports..."
	@env GOBIN=$(BIN_TOOLS_DIR) go install golang.org/x/tools/cmd/goimports@latest
