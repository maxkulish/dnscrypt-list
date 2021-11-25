ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
BIN_DIR = $(ROOT_DIR)/bin
PROJ_NAME = dnscrypt-list

help: _help_

_help_:
	@echo make build - build and push release with goreleaser. Output folder ./dist
	@echo make build-local - build local packages. Output folder ./dist

.PHONY: build
build:
	goreleaser build --rm-dist --snapshot -f .goreleaser.yml

build-local:
	goreleaser build --single-target --rm-dist --snapshot -f .goreleaser.yml