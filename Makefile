.PHONY: go_update

go_update:
	go version
	go get -u -v ./...
	go mod tidy
