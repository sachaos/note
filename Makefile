.PHONY: build

prepare:
	go get github.com/rakyll/statik
	go generate

build: prepare
	go build

install: prepare
	go install
