.PHONY: clean gen

clean:
	go clean

gen:
	cd ./cmd/gen && go run .

build-windows: gen
	CC=x86_64-w64-mingw32-gcc CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -o output/
