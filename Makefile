
PACKAGE := github.com/ShinyTrinkets/spinal

VERSION_VAR := main.VersionString
VERSION_VALUE ?= $(shell git describe --always --tags 2>/dev/null)
BUILD_TIME_VAR := main.BuildTime
BUILD_TIME_VALUE := $(shell date -u '+%Y-%m-%dT%I:%M:%S%z')

GOBUILD_LDFLAGS ?= \
	-X '$(VERSION_VAR)=$(VERSION_VALUE)' \
	-X '$(BUILD_TIME_VAR)=$(BUILD_TIME_VALUE)'

.PHONY: test clean deps build release

test:
	go vet -v
	go test -v ./parser

build: deps
	go install -x -ldflags "$(GOBUILD_LDFLAGS)"

deps:
	go get -x -ldflags "$(GOBUILD_LDFLAGS)"
	go get -t -x -ldflags "$(GOBUILD_LDFLAGS)"

release:
	GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w $(GOBUILD_LDFLAGS)" $(PACKAGE)
	mv spinal spinal-darwin
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w $(GOBUILD_LDFLAGS)" $(PACKAGE)
	mv spinal spinal-linux

clean:
	go clean -x -i ./...
