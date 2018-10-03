PROJECT ?= actionizer
VERSION ?= 1

# Current build settings
ARCH ?= amd64
NAME := $(PROJECT)-$(ARCH)
TARGET := $(shell pwd)/bin/$(NAME)
ENTRYPOINT := cmd/$(PROJECT).go

# Go command variables
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Linter variables
BIN_DIR := $(GOPATH)/bin
GOMETALINTER := $(BIN_DIR)/gometalinter

# Compiler variables
ALPINE_IMAGE = alpine:latest
GOLANG_IMAGE = library/golang:1.10
CVARS = CGO_ENABLED=0
CFLAGS = -a
LDFLAGS = -ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"
SRC ?= $(shell find . -type f -name '*.go' -not -path "./vendor/*")
PKGS = $(shell go list ./... | grep -v /vendor)

# Environment variables
include envfile
export $(shell sed 's/=.*//' envfile)

# Rules
#
.DEFAULT_GOAL: $(TARGET)
.PHONY: build clean run install test lint

$(TARGET): $(SRC)
	$(CVARS) $(GOBUILD) $(CFLAGS) $(LDFLAGS) -o $(TARGET) $(ENTRYPOINT)

build: $(TARGET)
	@true

clean:
	rm -f $(TARGET)

run: build
	@$(TARGET)

%:
	@true

install: 
	dep ensure -v

test: lint
	$(GOTEST) $(PKGS)

lint: $(GOMETALINTER)
	$(GOMETALINTER) ./... --vendor --fast --disable=maligned --enable=goimports

$(GOMETALINTER):
	$(GOGET) -u github.com/alecthomas/gometalinter
	$(GOMETALINTER) --install 1>/dev/null

fmt:
	goimports -l -w $(SRC)
