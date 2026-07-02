BINARY := fogos
CMD    := ./cmd/server

.PHONY: run build test test-go test-js fmt vet

run: fmt vet
	go run $(CMD)/...

build: fmt vet
	go build -o $(BINARY) $(CMD)/...

test: test-go test-js

test-go:
	go test ./...

test-js:
	node --test test/

fmt:
	gofmt -w .

vet:
	go vet ./...
