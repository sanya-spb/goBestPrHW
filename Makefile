PROJECT?=github.com/sanya-spb/goBestPrHW
PROJECTNAME=$(shell basename "$(PROJECT)")

GOOS?=linux
GOARCH?=amd64

CGO_ENABLED=1
EXE_FILE=dub_search

RELEASE := $(shell git tag -l | tail -1 | grep -E "v.+"|| echo devel)
COMMIT := git-$(shell git rev-parse --short HEAD)
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
COPYRIGHT := "sanya-spb"

## build: Build application
build:
	GOOS=${GOOS} GOARCH=${GOARCH} CGO_ENABLED=$(CGO_ENABLED) go build \
		-ldflags "-s -w \
		-X ${PROJECT}/pkg/version.version=${RELEASE} \
		-X ${PROJECT}/pkg/version.commit=${COMMIT} \
		-X ${PROJECT}/pkg/version.buildTime=${BUILD_TIME} \
		-X ${PROJECT}/pkg/version.copyright=${COPYRIGHT}" \
		-o ./cmd/dub_search/${EXE_FILE} ./cmd/dub_search/

## run: Run application
run:
	go run ./cmd/dub_search/

## clean: Clean build files
clean:
	go clean
	rm ./cmd/dub_search/$(EXE_FILE)

## linter: Run linters
linter:
	golangci-lint -c ./golangci-lint.yaml run

## help: Show this
help: Makefile
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
