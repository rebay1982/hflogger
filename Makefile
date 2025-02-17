.DEFAULT_GOAL := build

fmt:
	go fmt ./...

vet:
	go vet ./...

build: vet
	go build -v -o hflogger ./cmd/hflogger.go

test: build
	go test -v ./...
