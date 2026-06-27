BINARY := fogos
CMD    := ./cmd/server

.PHONY: run build test fmt vet

run: fmt vet
	go run $(CMD)/...

build: fmt vet
	go build -o $(BINARY) $(CMD)/...

test:
	go test ./...

fmt:
	gofmt -w .

vet:
	go vet ./...
