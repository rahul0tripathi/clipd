LINT_VERSION=1.41.1
LINT_IMAGE=golangci/golangci-lint:v${LINT_VERSION}-alpine
LINT_FLAGS=--timeout=10m0s

path :=$(if $(path), $(path), "./")

.PHONY: build-common
mod-setup:
	@ go version
	@ go clean
	@ go mod tidy && go mod download
	@ go mod verify

build-release: mod-setup
	@ CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o "./build/clipd" github.com/rahul0tripathi/clipd/cmd
