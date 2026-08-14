BINARY     := wa-persona-ai
CMD        := ./cmd/wa-persona-ai
BIN_DIR    := ./bin
BUILD_FLAGS := -ldflags="-s -w"

.PHONY: all build run clean test tidy lint docker-build docker-run

all: build

## build: compile the binary
build:
	@mkdir -p $(BIN_DIR)
	go build $(BUILD_FLAGS) -o $(BIN_DIR)/$(BINARY) $(CMD)
	@echo "Built: $(BIN_DIR)/$(BINARY)"

## run: build and run the bot
run: build
	$(BIN_DIR)/$(BINARY)

## dev: run without building (go run)
dev:
	go run $(CMD)

## test: run all tests
test:
	go test ./... -v -race

## tidy: tidy go modules
tidy:
	go mod tidy

## lint: run golangci-lint (requires golangci-lint installed)
lint:
	golangci-lint run ./...

## clean: remove build artifacts
clean:
	rm -rf $(BIN_DIR)

## docker-build: build Docker image
docker-build:
	docker build -t $(BINARY):latest .

## docker-run: run with docker-compose
docker-run:
	docker-compose up -d

## docker-stop: stop docker-compose
docker-stop:
	docker-compose down

## help: print this help
help:
	@grep -E '^## ' Makefile | sed 's/## /  /'
