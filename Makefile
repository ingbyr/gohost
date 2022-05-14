.PHONY: clean gen run build-linux build-windows

clean:
	-rm -r dal/query/*
	go clean

gen:
	cd ./cmd/gen && go run .

run:
	go run . -tags=memfs

build-linux: gen
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o output/

build-windows: gen
	CC=x86_64-w64-mingw32-gcc CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -o output/
