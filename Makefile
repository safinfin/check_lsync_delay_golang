BINNAME := check_lsync_delay

.PHONY: test
test:
	go test ./...

.PHONY: build
build:
	go build -o build/$(BINNAME) main.go

.PHONY: build-linux-amd64
build-linux-amd64:
	GOOS=linux GOARCH=amd64 go build -o build/linux-amd64/$(BINNAME) main.go

.PHONY: build-linux-arm64
build-linux-arm64:
	GOOS=linux GOARCH=arm64 go build -o build/linux-arm64/$(BINNAME) main.go
