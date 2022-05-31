ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
BIN_DIR = $(ROOT_DIR)/bin
PROJ_NAME = dnscrypt-list
SHELL := /bin/bash
FILES = $$(find . -type f -name "*.go" | grep -v '/vendor/')
VERSION ?= $(or $(shell git tag --sort=creatordate | grep -E '[0-9]' | tail -1 | cut -b 2-7 | awk -F. '{$$NF = $$NF + 1;} 1' | sed 's/ /./g'), $(shell echo 0.0.1))
PACKAGES = $$(go list ./... | grep -v '/vendor/')

help: _help_

_help_:
	@echo make fmt - fix formatting for the all files in the project
	@echo make version - show the current version of the project
	@echo make tag - create a new tag for the current commit and push it to the remote
	@echo make build - build and push release with goreleaser. Output folder ./dist
	@echo make build-local - build local package for current OS. Output folder ./dist
	@echo make test - run tests
	@echo make release-local - build and archive binaries for release. Output folder ./dist

.PHONY: build build-local test release-local vet fmt

build: .goreleaser.yml
	goreleaser build --rm-dist --snapshot -f .goreleaser.yml

build-local:
	goreleaser build --single-target --rm-dist --snapshot -f .goreleaser.yml

test:
	go test -timeout=15m -race -coverprofile=coverage.out -covermode=atomic -cover $(PACKAGES)

release-local:
	goreleaser release --snapshot --rm-dist -f .goreleaser.yml

vet:
	go vet ./...

fmt:
	gofmt -l -s -w $(FILES)

version:
	@echo $(VERSION)

coverage: test
	go tool cover -html=coverage.out