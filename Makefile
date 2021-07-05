PROJECT?=github.com/sanya-spb/goBestPrHW
PROJECTNAME=$(shell basename "$(PROJECT)")

GOOS?=linux
GOARCH?=amd64

CGO_ENABLED=1
EXE_FILE=csv-searcher

RELEASE := $(shell git tag -l | tail -1 | grep -E "v.+"|| echo devel)
COMMIT := git-$(shell git rev-parse --short HEAD)
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
COPYRIGHT := "sanya-spb"

## build: Build application
build:
	GOOS=${GOOS} GOARCH=${GOARCH} CGO_ENABLED=$(CGO_ENABLED) go build \
		-ldflags "-s -w \
		-X ${PROJECT}/pkg_tools/version.version=${RELEASE} \
		-X ${PROJECT}/pkg_tools/version.commit=${COMMIT} \
		-X ${PROJECT}/pkg_tools/version.buildTime=${BUILD_TIME} \
		-X ${PROJECT}/pkg_tools/version.copyright=${COPYRIGHT}" \
		-o ./cmd/csv-searcher/${EXE_FILE} ./cmd/csv-searcher/

## check: Run linters
check:
	golangci-lint -c ./golangci-lint.yaml run

## run: Run application
run:
	go run ./cmd/csv-searcher/

## clean: Clean build files
clean:
	go clean
	rm ./cmd/csv-searcher/${EXE_FILE}

## test: Run unit test
test:
	go test -v -short ${PROJECT}/cmd/csv-searcher/

## integration: Run integration test
integration:
	go test -v -run Integration ${PROJECT}/utils/fdouble/

## help: Show this
help: Makefile
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
