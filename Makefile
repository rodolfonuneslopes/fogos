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
	# Quoted so Node expands the glob itself: a bare directory argument is
	# resolved as a module entry point and fails with MODULE_NOT_FOUND.
	node --test 'test/**/*.test.js'

fmt:
	gofmt -w .

vet:
	go vet ./...
