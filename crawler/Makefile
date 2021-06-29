PROJECT?=github.com/sanya-spb/goBestPrHW
PROJECTNAME=$(shell basename "$(PROJECT)")

GOOS?=linux
GOARCH?=amd64

CGO_ENABLED=1
EXE_FILE=app_main

## build: Build application
build:
	GOOS=${GOOS} GOARCH=${GOARCH} CGO_ENABLED=$(CGO_ENABLED) go build \
		-o ${EXE_FILE} crawler.go main.go ugly_parser.go

## run: Run application
run: 
	go run .

## clean: Clean build files
clean:
	rm $(EXE_FILE)

## linter: Run linters
linter:
	golangci-lint -c ./golangci-lint.yaml run

## help: Show this
help: Makefile
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
