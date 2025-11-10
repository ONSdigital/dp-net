SHELL=bash

.PHONY: test
test:
	go test -v -count=1 -race -cover ./...

.PHONY: audit
audit: 
	dis-vulncheck 

.PHONY: build
build:
	go build ./...

.PHONY: lint
lint:
	golangci-lint run ./...
