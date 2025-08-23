include .env

LOCAL_BIN := "$(CURDIR)"/bin
YQ        := $(LOCAL_BIN)/yq
GOOSE 	  := $(LOCAL_BIN)/goose
LOCAL_MIGRATION_DSN := "$(shell DB_PASSWORD=$(DB_PASSWORD) $(YQ) -r '.postgres | "host=" + .host + " port=" + (.port|tostring) + " dbname=" + .db_name + " user=" + .user + " password=" + env(DB_PASSWORD)' "$(CURDIR)/config/local_config.yaml")"
LOCAL_MIGRATION_DIR := ./migrations

deps: ### deps tidy + verify
	go mod tidy && go mod verify
.PHONY: deps

format: ### run code formatter
	gofumpt -l -w .
	gci write . --skip-generated -s standard -s default
.PHONY: format

run: ### run app local
	go mod download && \
	go run cmd/app/main.go -config-path config/local_config.yaml -env-local
.PHONY: run

build: ### build app exe
	go mod download && \
	go build -o $(LOCAL_BIN)/app cmd/app/main.go
.PHONY: build

docker-build: build ### docker build app container
	docker build -t l0:v0.1 .
.PHONY: docker-build

compose-up: ### infrastructure up
	docker compose -f docker-compose.yaml up
.PHONY: compose-up

compose-down: ### infrastructure down 
	docker compose -f docker-compose.yaml down -v
.PHONY: compose-down

compose-with-app-up: ### infrastructure with app down 
	docker compose -f app.docker-compose.yaml up
.PHONY: compose-with-app-up

compose-with-app-down: ### infrastructure with app down 
	docker compose -f app.docker-compose.yaml down -v
.PHONY: compose-with-app-down

linter-golangci: ### check by golangci linter
	golangci-lint run
.PHONY: linter-golangci

linter-hadolint: ### check by hadolint linter
	git ls-files --exclude='Dockerfile*' --ignored | xargs hadolint
.PHONY: linter-hadolint

linter-dotenv: ### check by dotenv linter
	dotenv-linter
.PHONY: linter-dotenv

test: ### run test
	go test -v -race -covermode atomic ./internal/...
.PHONY: test

mock: ### run mockgen
	mockgen -source=./internal/repository/contracts.go -package=mock_repository -destination=./internal/repository/mocks/mock.go
	mockgen -source=./internal/usecase/contracts.go -package=usecase -destination=./internal/usecase/mocks/mock.go
.PHONY: mock

install-yq: ### install yq
	GOBIN=$(LOCAL_BIN) go install github.com/mikefarah/yq/v4@v4.45.1
.PHONY: install-yq

install-goose: ### install goose
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.24.2
.PHONY: install-goose

check-db-tools: ### check db tools
	@if [ ! -f $(LOCAL_BIN)/yq ]; then \
		echo "yq not found, installing..."; \
		$(MAKE) install-yq; \
	fi
	@if [ ! -f $(LOCAL_BIN)/goose ]; then \
		echo "goose not found, installing..."; \
		$(MAKE) install-goose; \
	fi
.PHONY: check-db-tools

migration-create: ### migration create
	$(MAKE) check-db-tools
	$(LOCAL_BIN)/goose -dir $(LOCAL_MIGRATION_DIR) create create_tables sql
.PHONY: migration-create

migration-status: ### migration status
	$(MAKE) check-db-tools
	$(LOCAL_BIN)/goose -dir $(LOCAL_MIGRATION_DIR) postgres $(LOCAL_MIGRATION_DSN) status -v
.PHONY: migration-status

migration-up: ### migration up
	$(MAKE) check-db-tools
	$(LOCAL_BIN)/goose -dir $(LOCAL_MIGRATION_DIR) postgres $(LOCAL_MIGRATION_DSN) up -v
.PHONY: migration-up

migration-down: ### migration down
	$(MAKE) check-db-tools
	$(LOCAL_BIN)/goose -dir $(LOCAL_MIGRATION_DIR) postgres $(LOCAL_MIGRATION_DSN) down -v
.PHONY: migration-down

pre-commit: mock format linter-golangci test ### run pre-commit
.PHONY: pre-commit
