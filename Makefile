.PHONY: mod

mod:
	go version
	go get -u -v ./...
	go mod tidy
