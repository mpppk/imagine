SHELL = /bin/bash

.PHONY: setup
setup:
	go get github.com/google/wire/cmd/wire
	go get github.com/goreleaser/goreleaser
	go get github.com/golang/mock/mockgen@v1.4.4
	go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.31.0

.PHONY: lint
lint: generate
	go vet ./...
	golangci-lint run
	goreleaser check

.PHONY: test
test: generate
	go test ./...

.PHONY: integration-test
integration-test:
	go test -tags=integration ./...

.PHONY: coverage
coverage: generate
	go test -race -coverprofile=coverage.txt -covermode=atomic ./...

.PHONY: codecov
codecov:  coverage
	bash <(curl -s https://codecov.io/bash)

.PHONY: wire
wire:
	go generate -tags=wireinject ./...

.PHONY: generate
generate: wire
	go generate ./...

.PHONY: build
build: generate
	go build

.PHONY: cross-build-snapshot
cross-build:
	goreleaser --rm-dist --snapshot

.PHONY: install
install:
	go install
