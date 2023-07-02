TARGET = bin
TARGET_BIN = app
GO_CMD_MAIN = cmd/app/main.go

MYSQL_DATABASE_URI = mysql://$(MYSQL_USERNAME):$(MYSQL_PASSWORD)@$(MYSQL_HOST):$(MYSQL_PORT)/$(MYSQL_DATABASE)

migrate:
	@echo \# make migrate name="$(name)"
	go run cmd/migrate/main.go create "$(name)"

migrate-up:
	atlas migrate apply \
  		--dir "file://migrations" \
  		--url $(MYSQL_DATABASE_URI)

generate-ent:
	go generate ./internal/ent

build:
	go build -o $(TARGET)/$(TARGET_BIN) $(GO_CMD_MAIN)

run:
	go run $(GO_CMD_MAIN)

run-test:
	go test --tags=integration ./api/server -count=1 -cover -race -v