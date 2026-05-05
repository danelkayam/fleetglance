APP_NAME := fleetglance

BIN_DIR := $(CURDIR)/bin
BUILD_DIR := $(CURDIR)/build
BIN_TOOLS_DIR := $(CURDIR)/bin-tools

BIN_AGENT := $(BIN_DIR)/agent
SRC_AGENT := $(CURDIR)/cmd/agent/main.go

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
	help build clean test lint format update  \
	build-deps clean-deps mod-download mod-tidy \
	tools-install tools-clean archive \
	build-agent build-docker-agent push-docker-agent

help:
	@echo "Available commands:"
	@echo "  make build         	 - Build the binary"
	@echo "  make build-agent   	 - Build the agent binary"
	@echo "  make clean         	 - Clean the build directory"
	@echo "  make test          	 - Run tests with coverage"
	@echo "  make lint          	 - Run linting"
	@echo "  make format        	 - Run code formatting"
	@echo "  make update        	 - Update go modules"
	@echo "  make mod-download  	 - Download go module dependencies"
	@echo "  make mod-tidy      	 - Tidy go module dependencies"
	@echo "  make tools-install 	 - Install development tools"
	@echo "  make tools-clean   	 - Clean development tools"
	@echo "  make archive       	 - Create code archive"
	@echo "  make build-docker-agent - Build agent docker image"
	@echo "  make push-docker-agent  - Push agent docker image to registry"

default: build

build: build-agent

build-agent:
	@echo "Building $(APP_NAME) agent..."
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
		go build -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/agent $(SRC_AGENT)

clean:
	@echo "Cleaning $(APP_NAME) build artifacts..."
	@rm -rf $(BIN_DIR)/*

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
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT=$(COMMIT) \
		--build-arg BUILT_AT=$(BUILT_AT) \
		-t $(AGENT_IMAGE):$(VERSION) \
		-t $(AGENT_IMAGE):latest

push-docker-agent:
	@echo "Pushing agent docker image..."
	@docker push $(AGENT_IMAGE):$(VERSION)
	@docker push $(AGENT_IMAGE):latest

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