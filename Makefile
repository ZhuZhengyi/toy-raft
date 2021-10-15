# Makefile
#
#
#

all:
	@echo "all done"

build:
	@go build -o build/toy-raft

test:
	@echo "test"
	@go test ./...
