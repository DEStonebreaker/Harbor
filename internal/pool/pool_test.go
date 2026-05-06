package pool

import (
	"sync"
	"testing"
)

func TestNewPool_Valid(t *testing.T) {
	urls := []string{"http://a:8080", "http://b:8080", "http://c:8080"}
	p, err := NewPool(urls)
	if err != nil {
		t.Fatalf("NewPool: %v", err)
	}
	if len(p.Backends) != len(urls) {
		t.Errorf("want %d backends, got %d", len(urls), len(p.Backends))
	}
	for _, b := range p.Backends {
		if !b.Alive {
			t.Errorf("backend %s should start Alive=true", b.URL)
		}
		if b.URL == nil {
			t.Errorf("backend URL nil")
		}
	}
}

func TestNewPool_BadURL(t *testing.T) {
	// url.Parse is lenient. Use unclosed bracket which it rejects.
	_, err := NewPool([]string{"http://[::1"})
	if err == nil {
		t.Fatal("expected error for malformed URL")
	}
}

func TestNewPool_Empty(t *testing.T) {
	p, err := NewPool(nil)
	if err != nil {
		t.Fatalf("NewPool(nil): %v", err)
	}
	if len(p.Backends) != 0 {
		t.Errorf("want 0 backends, got %d", len(p.Backends))
	}
	if p.Alive() == nil {
		t.Error("Alive() should return empty slice not nil")
	}
}

func TestPool_Alive_All(t *testing.T) {
	p, _ := NewPool([]string{"http://a", "http://b"})
	got := p.Alive()
	if len(got) != 2 {
		t.Errorf("want 2 alive, got %d", len(got))
	}
}

func TestPool_Alive_None(t *testing.T) {
	p, _ := NewPool([]string{"http://a", "http://b"})
	for _, b := range p.Backends {
		p.SetAlive(b, false)
	}
	got := p.Alive()
	if len(got) != 0 {
		t.Errorf("want 0 alive, got %d", len(got))
	}
	if got == nil {
		t.Error("expected empty slice, got nil")
	}
}

func TestPool_Alive_Mixed(t *testing.T) {
	p, _ := NewPool([]string{"http://a", "http://b", "http://c"})
	p.SetAlive(p.Backends[1], false)

	got := p.Alive()
	if len(got) != 2 {
		t.Fatalf("want 2 alive, got %d", len(got))
	}
	for _, b := range got {
		if b == p.Backends[1] {
			t.Error("dead backend leaked into Alive()")
		}
	}
}

func TestPool_SetAlive_Toggle(t *testing.T) {
	p, _ := NewPool([]string{"http://a"})
	b := p.Backends[0]

	p.SetAlive(b, false)
	if b.Alive {
		t.Error("SetAlive(false) did not flip flag")
	}

	p.SetAlive(b, true)
	if !b.Alive {
		t.Error("SetAlive(true) did not flip flag")
	}
}

// Run with: go test -race ./internal/pool/
// Without -race, conflicts go undetected.
func TestPool_RaceSafe(t *testing.T) {
	p, _ := NewPool([]string{"http://a", "http://b", "http://c"})

	const iterations = 500
	var wg sync.WaitGroup

	for i := 0; i < iterations; i++ {
		wg.Add(3)
		go func() {
			defer wg.Done()
			_ = p.Alive()
		}()
		go func(idx int) {
			defer wg.Done()
			p.SetAlive(p.Backends[idx%len(p.Backends)], false)
		}(i)
		go func(idx int) {
			defer wg.Done()
			p.SetAlive(p.Backends[idx%len(p.Backends)], true)
		}(i)
	}
	wg.Wait()
}
