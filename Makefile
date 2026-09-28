.PHONY: build run test vet run-hurl stop dev dev-down

HURL_DIR := tools/_hurl
BASE_URL ?= http://localhost:8080

PODMAN_DOCKER_HOST := $(shell tools/podman.sh 2>/dev/null || true)

build: # [ make build ]
	go build -o bin/peeko ./cmd/peeko
	@echo "> Peeko Built"

run: # [ make run ]
	@test -f .env || cp .env.sample .env
	@command -v air >/dev/null 2>&1 || { echo "> air not found — install: go install github.com/air-verse/air@latest"; exit 1; }
	@set -a && . ./.env && set +a && echo "> UI: http://localhost:$${PORT:-8080}/ui" && air

stop: # [ make stop ]
	@lsof -ti :8080 -s TCP:LISTEN | xargs kill 2>/dev/null || true
	@echo "> Peeko Stopped"

vet: # [ make vet ]
	go vet ./...

test: # [ make test ]
	go test ./...

run-hurl: # [ make run-hurl BASE_URL=http://localhost:8080 ]
	bash tools/run-hurl.sh $(BASE_URL)

dev: # [ make dev ]
	@test -f .env || cp .env.sample .env
	$(PODMAN_DOCKER_HOST) docker compose -f dev.docker-compose.yaml up --build -d
	@echo "> UI: http://localhost:8080/ui"

dev-down: # [ make dev-down ]
	$(PODMAN_DOCKER_HOST) docker compose -f dev.docker-compose.yaml down
