
@PHONY: dev build clean

.DEFAULT_GOAL := dev
CGO_ENABLED := 0

dev:
	go build -o dist/ethcli

build: dist/ethcli-linux dist/ethcli-mac

dist/ethcli-linux:
	GOOS=linux  GOARCH=amd64 go build -o dist/ethcli-linux

dist/ethcli-mac:
	GOOS=darwin GOARCH=amd64 go build -o dist/ethcli-darwin-amd64
	GOOS=darwin GOARCH=arm64 go build -o dist/ethcli-darwin-arm64
	lipo -create -output dist/ethcli-mac dist/ethcli-darwin*
	rm dist/ethcli-darwin*

clean:
	rm -rf dist/
