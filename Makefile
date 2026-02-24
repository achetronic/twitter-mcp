.PHONY: build run clean tidy

BINARY_NAME=twitter-mcp
BUILD_DIR=bin

build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/main.go

run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BUILD_DIR)/$(BINARY_NAME) -config config.yaml

tidy:
	@echo "Tidying dependencies..."
	go mod tidy

clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)

install: build
	@echo "Installing $(BINARY_NAME)..."
	cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
