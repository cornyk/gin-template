# Makefile
APP_API_NAME := app-api
APP_CLI_NAME := app-cli
BUILD_TARGET_DIR := build/target

.PHONY: build-test clean

build-test:
	make clean
	mkdir -p $(BUILD_TARGET_DIR)/runtime/logs
	go build -o $(BUILD_TARGET_DIR)/$(APP_API_NAME) ./cmd/api/main.go
	go build -o $(BUILD_TARGET_DIR)/$(APP_CLI_NAME) ./cmd/cli/main.go
	cp build/config/config.test.yaml $(BUILD_TARGET_DIR)/config.yaml

clean:
	rm -rf $(BUILD_TARGET_DIR)/*
