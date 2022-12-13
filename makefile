.PHONY: linux windows macos

linux:
	GOOS=linux GOARCH=amd64 go build -o build/linux/gohost

windows:
	GOOS=windows GOARCH=amd64 go build -o build/win/gohost.exe
	cp gohost.bat build/win/

macos:
