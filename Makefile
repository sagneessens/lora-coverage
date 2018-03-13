.PHONY: all build clean
VERSION := 0.1
COMMIT := $(shell git describe --always)
GOOS ?= darwin
GOARCH ?= amd64
GOPATH ?= $(HOME)/go/
BUILD_DATE = `date -u +%Y-%m-%dT%H:%M.%SZ`

all: clean build

build:
	@echo "Compiling source for $(GOOS) $(GOARCH)"
	@mkdir -p build
	@GOOS=$(GOOS) GOARCH=$(GOARCH) go build -a -ldflags "-X main.version=$(VERSION) -X main.build=$(COMMIT) -X main.buildDate=$(BUILD_DATE)" -o build/lora-coverage-$(GOOS)-$(GOARCH)$(BINEXT) main.go

clean:
	@echo "Cleaning up workspace"
	@rm -rf build
	@rm -rf coverage.log
