ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
BIN_DIR = $(ROOT_DIR)/bin
PROJ_NAME = dnscrypt-list
SHELL := /bin/bash

help: _help_

_help_:
	@echo make build - build and push release with goreleaser. Output folder ./dist
	@echo make build-local - build local package for current OS. Output folder ./dist
	@echo make test - run tests
	@echo make release-local - build and archive binaries for release. Output folder ./dist

.PHONY: build
build: .goreleaser.yml
	goreleaser build --rm-dist --snapshot -f .goreleaser.yml

.PHONY: build-local
build-local:
	goreleaser build --single-target --rm-dist --snapshot -f .goreleaser.yml

.PHONY: test
test:
	go test -v ./...

.PHONY: release-local
release-local:
	goreleaser release --snapshot --rm-dist -f .goreleaser.yml