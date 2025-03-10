dev:
	go run ./cmd/gateway/gateway.go

test:
	go test ./... -v --cover

build:
	go build -o bin/gateway ./cmd/gateway

clean:
	rm -rf bin
