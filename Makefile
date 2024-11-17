# Variables
APP_NAME := piscator
BUILD_DIR := ./bin
SRC_DIR := ./cmd/$(APP_NAME)
VERSION := $(shell git describe --tags --always --dirty)
LDFLAGS := -s -w -X main.version=$(VERSION)

# Targets
.PHONY: build debug test tidy clean

clean:
	rm -rf ${BUILD_DIR}

tidy: clean
	go mod tidy
	go mod vendor
	go fmt ./...

test:
	staticcheck ./...
	go test -v -race -buildvcs ./...

build: tidy
	go build -o $(BUILD_DIR)/$(APP_NAME) -ldflags "-s -w $(LDFLAGS)" $(SRC_DIR)

debug: tidy
	go build -o $(BUILD_DIR)/$(APP_NAME)-debug -ldflags "$(LDFLAGS)" $(SRC_DIR)
