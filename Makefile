BIN_DIR := bin
LINUX_BIN_DIR := $(BIN_DIR)/linux
MACOS_BIN_DIR := $(BIN_DIR)/darwin

BINARY_NAME := ddlogexporter

GOARCH ?= amd64

build-linux:
	mkdir -p $(LINUX_BIN_DIR)
	GOOS=linux GOARCH=$(GOARCH) go build -o $(LINUX_BIN_DIR)/$(BINARY_NAME) cmd/main.go

build-macos:
	mkdir -p $(MACOS_BIN_DIR)
	GOOS=darwin GOARCH=$(GOARCH) go build -o $(MACOS_BIN_DIR)/$(BINARY_NAME) cmd/main.go

all: build-linux build-macos
