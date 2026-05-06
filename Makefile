.PHONY: run build test test-race test-integration cover clean

BIN := bin/harbor
PKG := ./...

run:
	go run ./backend

build:
	go build -o $(BIN) ./backend

# Fast unit tests. Run on every save.
test:
	go test $(PKG)

# Detect data races. Slower but catches concurrency bugs.
test-race:
	go test -race $(PKG)

# Slow tests behind build tag. Real network, real proxy.
test-integration:
	go test -tags=integration ./test/...

# HTML coverage report.
cover:
	go test -coverprofile=coverage.out $(PKG)
	go tool cover -html=coverage.out -o coverage.html
	@echo "open coverage.html"

clean:
	rm -rf bin coverage.out coverage.html
