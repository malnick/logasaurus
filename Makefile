VERSION := $(shell git describe --tags)
REVISION := $(shell git rev-parse --short HEAD)

BINARY_NAME := loga 

LDFLAGS := -X github.com/malnick/logasaurus/config.VERSION=$(VERSION) -X github.com/malnick/logasaurus/config.REVISION=$(REVISION) 

FILES := $(shell go list ./... | grep -v vendor)

all: test install

test:
	@echo "+$@"
	go test $(FILES)  -cover

build: 
	@echo "+$@"
	go build -v -o loga_$(VERSION) -ldflags '$(LDFLAGS)' main.go

install:
	@echo "+$@"
	go install -v -ldflags '$(LDFLAGS)' $(FILES)
