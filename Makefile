.PHONY: clean gen

clean:
	go clean

gen:
	cd ./cmd/gen && go run .