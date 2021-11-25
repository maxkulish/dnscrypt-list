ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
BIN_DIR = $(ROOT_DIR)/bin
PROJ_NAME = dnscrypt-list
SHELL := /bin/bash

help: _help_

_help_:
	@echo make build - build and push release with goreleaser. Output folder ./dist
	@echo make build-local - build local package for current OS. Output folder ./dist
	@echo make test - run tests

.PHONY: build
build: .goreleaser.yml
	goreleaser build --rm-dist --snapshot -f .goreleaser.yml

build-local:
	goreleaser build --single-target --rm-dist --snapshot -f .goreleaser.yml

test:
	go test -v ./...