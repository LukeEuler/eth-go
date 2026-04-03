all: build

CILINT := $(shell command -v golangci-lint 2> /dev/null)

format:
ifndef CILINT
	$(error "golangci-lint is not available please install golangci-lint")
endif
	golangci-lint fmt

cilint:
ifndef CILINT
	$(error "golangci-lint is not available please install golangci-lint")
endif
	golangci-lint run --timeout 5m0s

test: cilint
	go test -cover ./...

build: test
	go build -o build/main app/main.go

.PHONY: format test build
