# Makefile
#
#
#

BUILD_TARGETS := build/toy-raft

PHONY := all
all: generate build
	@echo "make all done"

PHONY += generate
generate:
	@go generate ./...

PHONY += build
build: $(BUILD_TARGETS)
	@go build -o build/toy-raft

PHONY += clean
clean:
	@rm -rf build/* raft/*_string.go

PHONY += test
test: generate
	@go test -v ./...

PHONY += docs
docs:
	@godoc -http=:6060

.phony: $(PHONY)
