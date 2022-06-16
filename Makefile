
VERSION_VAR := main.Version
VERSION_VALUE ?= $(shell cat VERSION.txt)
REVISION_VAR := main.CommitHash
REVISION_VALUE ?= $(shell git rev-parse --short HEAD 2>/dev/null)
BUILT_VAR := main.BuildTime
BUILT_VALUE := $(shell date -u '+%Y-%m-%dT%I:%M:%S%z')

GOBUILD_LDFLAGS ?= \
	-X '$(VERSION_VAR)=$(VERSION_VALUE)' \
	-X '$(REVISION_VAR)=$(REVISION_VALUE)' \
	-X '$(BUILT_VAR)=$(BUILT_VALUE)'

# Option for version bump
BUMP ?= patch

.PHONY: test clean build release version

test:
	go test -v ./state
	go test -v ./parser

build:
	go build -o spin -x -ldflags "$(GOBUILD_LDFLAGS)"
	mv spin ~/go/bin/spin

release:
	GOOS=darwin GOARCH=amd64 go build -o spin-darwin -ldflags "-s -w $(GOBUILD_LDFLAGS)"
	GOOS=linux GOARCH=amd64 go build -o spin-linux -ldflags "-s -w $(GOBUILD_LDFLAGS)"

version:
	python version_bump.py `cat VERSION.txt` --$(BUMP) > VERSION.txt

clean:
	go clean -x -i -cache -testcache ./...
