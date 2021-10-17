# Makefile
#
#
#

BUILD_TARGETS := build/toy-raft

PHONY := all
all: test build
	@echo "make all done"

PHONY += build
build: $(BUILD_TARGETS)
	@go generate ./...
	@go build -o build/toy-raft

PHONY += clear
clear:
	@rm -rf build/*

PHONY += test
test:
	@go test ./...

.phony: $(PHONY)
