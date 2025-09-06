.PHONY: mod

mod:
	go get -u -v ./...
	go mod tidy
