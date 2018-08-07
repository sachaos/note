.PHONY: build

prepare:
	dep ensure
	(cd assets && npm install)
	go get github.com/rakyll/statik
	go generate

build: prepare
	go build

install: prepare
	go install
