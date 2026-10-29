SHELL := /bin/bash
ENV_FILE ?= .env
GO ?= go
BIN_DIR := bin

ENV_LOAD = set -a; if [ -f "$(ENV_FILE)" ]; then . "./$(ENV_FILE)"; fi; set +a;

.PHONY: help build run-web run-worker health metrics test test-race vet fmt fmt-check check tidy versions clean

help:
	@echo "Доступные цели:"
	@echo "  make build        собрать bin/web и bin/worker"
	@echo "  make run-web      запустить web с переменными из $(ENV_FILE)"
	@echo "  make run-worker   запустить worker с переменными из $(ENV_FILE)"
	@echo "  make health       GET /healthz"
	@echo "  make metrics      GET /metrics на внутреннем порту"
	@echo "  make test         go test ./..."
	@echo "  make test-race    go test -race ./..."
	@echo "  make vet          go vet ./..."
	@echo "  make fmt          gofmt -w cmd internal"
	@echo "  make fmt-check    проверить форматирование без изменений"
	@echo "  make check        fmt-check + vet + test + build"
	@echo "  make tidy         go mod tidy"
	@echo "  make versions     версии Go, templ, sqlc"
	@echo "  make clean        удалить $(BIN_DIR)"

build:
	@mkdir -p $(BIN_DIR)
	$(GO) build -o $(BIN_DIR)/web ./cmd/web
	$(GO) build -o $(BIN_DIR)/worker ./cmd/worker

run-web:
	@$(ENV_LOAD) $(GO) run ./cmd/web

run-worker:
	@$(ENV_LOAD) $(GO) run ./cmd/worker

health:
	@$(ENV_LOAD) curl -fsS "$${APP_BASE_URL:-http://localhost:8080}/healthz"; echo

metrics:
	@$(ENV_LOAD) curl -fsS "http://$${METRICS_ADDR:-127.0.0.1:9090}/metrics" | head -20

test:
	$(GO) test ./...

test-race:
	$(GO) test -race ./...

vet:
	$(GO) vet ./...

fmt:
	gofmt -w cmd internal

fmt-check:
	@out="$$(gofmt -l cmd internal)"; \
	if [ -n "$$out" ]; then echo "Не отформатировано:"; echo "$$out"; exit 1; fi

check: fmt-check vet test build

tidy:
	$(GO) mod tidy

versions:
	@$(GO) version
	@$(GO) tool templ version
	@$(GO) tool sqlc version

clean:
	rm -rf $(BIN_DIR)
