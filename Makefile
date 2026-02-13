WIN_GO_VERSION=1.20.14
WIN_GO_BIN=$(HOME)/sdk/go$(WIN_GO_VERSION)/bin/go

# Windows build (with icon)
build-windows:
	powershell -ExecutionPolicy Bypass -File scripts/build-windows.ps1 -Hidden -OutputName "蟹掌柜.exe"

build-windows-legacy:
	$(WIN_GO_BIN) build -C src -ldflags="-s -w -H=windowsgui" -o ../dist/weblauncher.exe

build-mac:
	go build -C src -ldflags="-s -w" -o ../dist/weblauncher-mac

build-linux:
	GOOS=linux GOARCH=amd64 go build -C src -ldflags="-s -w" -o ../dist/weblauncher-linux

run:
	go run -C src

# Generate syso manually
gen-syso:
	cd src && rsrc -ico="assets/icon.ico" -o="rsrc.syso"

# Clean build files
clean:
	rm -f dist/* src/*.syso
