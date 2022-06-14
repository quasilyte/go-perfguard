GOPATH_DIR=`go env GOPATH`
VERSION=$(shell git describe --tags)

test:
	go test -count 2 -coverpkg=./... -coverprofile=coverage.txt -covermode=atomic ./...
	go test -bench=. ./...
	@echo "everything is OK"

ci-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH_DIR)/bin v1.43.0
	$(GOPATH_DIR)/bin/golangci-lint run ./...
	go install github.com/quasilyte/go-consistent@latest
	$(GOPATH_DIR)/bin/go-consistent ./cmd/... ./perfguard/... ./internal/...
	go build -o bin/perfguard ./cmd/perfguard && ./bin/perfguard lint ./...
	go run ./_script/check.go
	@echo "everything is OK"

lint:
	golangci-lint run ./...
	go run ./_script/check.go
	@echo "everything is OK"

.PHONY: ci-lint lint test

build:
	go build -o bin/perfguard -ldflags "-s -w -X ./cmd/perfguard.BuildVersion=${VERSION}" ./cmd/perfguard

draft-release:
	go run releaser/release.go -version ${VERSION}
