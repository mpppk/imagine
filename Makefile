SHELL = /bin/bash

.PHONY: setup
setup:
	go get github.com/google/wire/cmd/wire
	go get github.com/goreleaser/goreleaser
	go get github.com/rakyll/statik
	go get github.com/golang/mock/mockgen@v1.4.4
	git checkout go.mod
	git checkout go.sum

.PHONY: clean
clean:
	rm -rf static/out
	rm -rf statik
	rm -f imagine

.PHONY: lint
lint: generate
	go vet ./...
	goreleaser check

.PHONY: lint-all
lint-all: lint
	golangci-lint run

.PHONY: test
test: lint
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

.PHONY: go-generate
go-generate: wire
	go generate ./...

.PHONY: generate
generate: go-generate
	yarn --cwd static export
	statik -f -src static/out

.PHONY: build
build: test
	go build

.PHONY: cross-build-snapshot
cross-build-snapshot: test
	goreleaser --rm-dist --snapshot

.PHONY: install
install: build
	go install
