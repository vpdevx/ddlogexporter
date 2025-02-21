BIN_DIR := bin
LINUX_BIN_DIR := $(BIN_DIR)/linux
MACOS_BIN_DIR := $(BIN_DIR)/darwin
WIN_BIN_DIR := $(BIN_DIR)/windows

BINARY_NAME := ddlogexporter

GOARCH ?= amd64

build-linux:
	mkdir -p $(LINUX_BIN_DIR)
	GOOS=linux GOARCH=$(GOARCH) go build -o $(LINUX_BIN_DIR)/$(BINARY_NAME) cmd/main.go

build-macos:
	mkdir -p $(MACOS_BIN_DIR)
	GOOS=darwin GOARCH=$(GOARCH) go build -o $(MACOS_BIN_DIR)/$(BINARY_NAME) cmd/main.go

build-windows:
	mkdir -p $(WIN_BIN_DIR)
	GOOS=windows GOARCH=$(GOARCH) go build -o $(WIN_BIN_DIR)/$(BINARY_NAME).exe cmd/main.go

all: build-linux build-macos build-windows
