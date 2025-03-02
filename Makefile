# Define variables
BINARY_NAME=gollama
BUILD_DIR=build
SRC_DIR=src

# Default target
all: build

# Build the binary
build:
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(SRC_DIR)/main.go

# Clean the build directory
clean:
	@rm -rf ./db
	@rm -rf $(BUILD_DIR)

# Phony targets
.PHONY: all build clean
