# Makefile for kubeRadar

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
BINARY_NAME=kubeRadar
# OS detection for binary extension and mkdir command
ifeq ($(OS),Windows_NT)
	BINARY_NAME := kubeRadar.exe
	MKDIR = if not exist $(BUILD_DIR) mkdir $(BUILD_DIR)
else
	BINARY_NAME := kubeRadar
	MKDIR = mkdir -p $(BUILD_DIR)
endif
BINARY_UNIX=$(BINARY_NAME)_unix

# Build parameters
BUILD_DIR=build
MAIN_PATH=main.go

.PHONY: all build clean test deps

all: test build

build:
	@$(MKDIR)
	@$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

build-linux:
	@$(MKDIR)
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_UNIX) $(MAIN_PATH)

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BUILD_DIR)/$(BINARY_NAME)
	rm -f $(BUILD_DIR)/$(BINARY_UNIX)

deps:
	$(GOCMD) mod download

run:
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	./$(BUILD_DIR)/$(BINARY_NAME)
