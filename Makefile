VERSION=0.0.11
LDFLAGS=-ldflags "-w -s -X main.version=${VERSION}"

all: mackerel-plugin-pinging

.PHONY: mackerel-plugin-pinging linux check lint

mackerel-plugin-pinging: *.go
	go build $(LDFLAGS) -o mackerel-plugin-pinging

linux: *.go
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o mackerel-plugin-pinging

check:
	go test ./...

lint:
	golangci-lint run ./...