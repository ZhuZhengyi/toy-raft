# Makefile
#
#
#

BUILD_TARGETS := build/toy-raft

PHONY := all
all: generate test build
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
test:
	@go test ./...

PHONY += docs
docs:
	@godoc -http=:6060

.phony: $(PHONY)
