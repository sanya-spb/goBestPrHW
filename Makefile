PROJECT?=github.com/sanya-spb/goBestPrHW
PROJECTNAME=$(shell basename "$(PROJECT)")

GOOS?=linux
GOARCH?=amd64

CGO_ENABLED=1
EXE_FILE=app_main

RELEASE := $(shell git tag -l | tail -1 | grep -E "v.+"|| echo devel)
COMMIT := git-$(shell git rev-parse --short HEAD)
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
COPYRIGHT := "sanya-spb"

## build: Build application
build:
	GOOS=${GOOS} GOARCH=${GOARCH} CGO_ENABLED=$(CGO_ENABLED) go build \
		-ldflags "-s -w -X ${PROJECT}/utils/version.version=${RELEASE} \
		-X ${PROJECT}/utils/version.commit=${COMMIT} \
		-X ${PROJECT}/utils/version.buildTime=${BUILD_TIME} \
		-X ${PROJECT}/utils/version.copyright=${COPYRIGHT}" \
		-o ${EXE_FILE} main.go

## run: Run application
run: 
	go run .

## clean: Clean build files
clean:
	go clean
	rm $(EXE_FILE)

## test: Run autotest
test:
	# go test -v ${PROJECT}/utils/config/
	go test -v

## help: Show this
help: Makefile
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
