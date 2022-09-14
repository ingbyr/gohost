.PHONY: linux windows macos

linux:
	GOOS=linux GOARCH=amd64 go build -o build/gohost

windows:
	GOOS=windows GOARCH=amd64 go build -o build/gohost.exe

macos:
