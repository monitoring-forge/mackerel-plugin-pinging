VERSION=0.0.8
GITCOMMIT?=$(shell git describe --dirty --always 2>/dev/null)
LDFLAGS=-ldflags "-w -s -X main.version=${VERSION} -X main.commit=${GITCOMMIT}"

all: mackerel-plugin-pinging

.PHONY: mackerel-plugin-pinging

mackerel-plugin-pinging: main.go
	go build $(LDFLAGS) -o mackerel-plugin-pinging

linux: main.go
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o mackerel-plugin-pinging

check:
	go test ./...
