build:
	go build -o bin/cockpit ./cmd/tui

run: build
	./bin/cockpit

test:
	go test ./...
