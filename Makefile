.DEFAULT_GOAL := build
.PHONY: build

VERSION="v0.4"
GOOS="linux"
GOARCH="amd64"

build:
	@CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o whale-$(GOOS)-$(GOARCH)-$(VERSION) cmd/main.go