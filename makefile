.PHONY: build
build:
	go build -v ./cmd/main.go

.PHONY: test
test:
	go test -v -race -timeout 30s ./...

.PHONY: cover
cover:
	go test -coverprofile cover.out ./... && go tool cover -html=cover.out

.PHONY: golines
golines:
	~/go/bin/golines -m 88 -w .

.DEFAULT_GOAL := build