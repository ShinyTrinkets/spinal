
PACKAGE := github.com/ShinyTrinkets/gears.go

VERSION_VAR := $(PACKAGE).VersionString
VERSION_VALUE ?= $(shell git describe --always --tags 2>/dev/null)
REV_VAR := $(PACKAGE).RevisionString
REV_VALUE ?= $(shell git rev-parse HEAD --short 2>/dev/null)
BUILD_TIME_VAR := $(PACKAGE).BuildTime
BUILD_TIME_VALUE := $(shell date -u '+%Y-%m-%d-%I:%M:%S-%Z')

GOBUILD_LDFLAGS ?= \
	-X '$(VERSION_VAR)=$(VERSION_VALUE)' \
	-X '$(REV_VAR)=$(REV_VALUE)' \
	-X '$(BUILD_TIME_VAR)=$(BUILD_TIME_VALUE)'

.PHONY: test
test:
	go test -v ./parser

.PHONY: build
build: deps
	go install -x -ldflags "$(GOBUILD_LDFLAGS)"

.PHONY: deps
deps:
	go get -x -ldflags "$(GOBUILD_LDFLAGS)"
	go get -t -x -ldflags "$(GOBUILD_LDFLAGS)"

.PHONY: clean
clean:
	go clean -i ./...
