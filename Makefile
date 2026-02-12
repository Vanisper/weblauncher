# Makefile
WIN_GO_VERSION=1.20.14
WIN_GO_BIN=$(HOME)/sdk/go$(WIN_GO_VERSION)/bin/go

build-win7:
	$(WIN_GO_BIN) build -ldflags="-s -w -H=windowsgui" -o dist/weblauncher-win7.exe

build-mac:
	go build -ldflags="-s -w" -o dist/weblauncher-mac

build-linux:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o dist/weblauncher-linux
