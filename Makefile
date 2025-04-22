# Configurações
APP_NAME := app
CMD_PATH := ./cmd/app
BIN_DIR := ./bin
VERSION := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
COMMIT := $(shell git rev-parse --short HEAD)
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# ===== Targets =====

## Build local
build:
	@echo "🔨 Buildando localmente..."
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
		-ldflags="-s -w -X 'main.version=$(VERSION)' -X 'main.commit=$(COMMIT)' -X 'main.date=$(BUILD_DATE)'" \
		-o $(BIN_DIR)/$(APP_NAME) $(CMD_PATH)
	@echo "✅ Binário gerado em $(BIN_DIR)/$(APP_NAME)"

## Build com Docker
docker-build:
	@echo "🐳 Buildando com Docker..."
	docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT=$(COMMIT) \
		--build-arg BUILD_DATE=$(BUILD_DATE) \
		-f Dockerfile.build \
		-o type=local,dest=$(BIN_DIR) .

## Limpa binários
clean:
	rm -rf $(BIN_DIR)
	@echo "🧹 Limpeza feita!"

.PHONY: build docker-build clean
