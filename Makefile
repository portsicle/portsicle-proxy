APP_NAME = attorney-toolkit
PKG_DIRS = ./cmd ./internal
OUTPUT_DIR = ./bin
GO_FILES = $(shell find $(PKG_DIRS) -type f -name '*.go')

.PHONY: all
all: build

# Build the main application
.PHONY: build
build: $(GO_FILES)
	@echo "Building $(APP_NAME)..."
	go build -o $(OUTPUT_DIR)/$(APP_NAME) ./main.go
	@echo "Build complete: $(OUTPUT_DIR)/$(APP_NAME)"

# Run the application
.PHONY: run
run:
	@echo "Running $(APP_NAME)..."
	go run ./main.go

# Format the code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...
