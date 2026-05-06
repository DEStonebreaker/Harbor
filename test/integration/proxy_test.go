//go:build integration

// Integration tests. Spin up real backends + Harbor proxy, hit over HTTP.
// Run: go test -tags=integration ./test/...
// Skipped by default `go test ./...` so unit tests stay fast.

package integration

import "testing"

func TestPlaceholder(t *testing.T) {
	t.Skip("integration scaffold — wire up real proxy + backends here once handler exists")
}
