package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/rodolfonuneslopes/fogos/internal/fogos"
	"github.com/rodolfonuneslopes/fogos/internal/handler"
)

// countingClient records how many times ActiveIncidents was called, per
// concelho, so tests can assert on cache/singleflight behavior.
type countingClient struct {
	mu     sync.Mutex
	calls  map[string]int
	result []fogos.Incident
	err    error
}

func newCountingClient() *countingClient {
	return &countingClient{calls: make(map[string]int)}
}

func (c *countingClient) ActiveIncidents(concelho string) ([]fogos.Incident, error) {
	c.mu.Lock()
	c.calls[concelho]++
	c.mu.Unlock()
	if c.err != nil {
		return nil, c.err
	}
	return c.result, nil
}

func (c *countingClient) callCount(concelho string) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.calls[concelho]
}

// blockingClient blocks every call until release is closed, so concurrent
// requests can be forced to overlap and exercise singleflight collapsing.
type blockingClient struct {
	release chan struct{}
	calls   int32
}

func (b *blockingClient) ActiveIncidents(concelho string) ([]fogos.Incident, error) {
	atomic.AddInt32(&b.calls, 1)
	<-b.release
	return []fogos.Incident{{ID: "1", Concelho: concelho}}, nil
}

func TestIncidents_UnknownConcelho(t *testing.T) {
	mux := handler.NewMux(newCountingClient())
	req := httptest.NewRequest(http.MethodGet, "/api/incidents?concelho=Nowhereville", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", rec.Code)
	}
}

func TestIncidents_KnownConcelho(t *testing.T) {
	client := newCountingClient()
	client.result = []fogos.Incident{{ID: "1", Concelho: "Lisboa"}}
	mux := handler.NewMux(client)

	req := httptest.NewRequest(http.MethodGet, "/api/incidents?concelho=Lisboa", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("want 200, got %d: %s", rec.Code, rec.Body.String())
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("want application/json, got %q", ct)
	}
	if !strings.Contains(rec.Body.String(), "Lisboa") {
		t.Errorf("body missing expected incident: %s", rec.Body.String())
	}
}

func TestIncidents_UpstreamError(t *testing.T) {
	client := newCountingClient()
	client.err = errUpstream
	mux := handler.NewMux(client)

	req := httptest.NewRequest(http.MethodGet, "/api/incidents", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadGateway {
		t.Fatalf("want 502, got %d", rec.Code)
	}
}

func TestIncidents_CacheAvoidsSecondCall(t *testing.T) {
	client := newCountingClient()
	client.result = []fogos.Incident{{ID: "1", Concelho: "Porto"}}
	mux := handler.NewMux(client)

	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodGet, "/api/incidents?concelho=Porto", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("request %d: want 200, got %d", i, rec.Code)
		}
	}

	if got := client.callCount("Porto"); got != 1 {
		t.Errorf("want 1 upstream call across 2 cached requests, got %d", got)
	}
}

func TestIncidents_SingleflightCollapsesConcurrentMisses(t *testing.T) {
	client := &blockingClient{release: make(chan struct{})}
	mux := handler.NewMux(client)

	const concurrent = 5
	var wg sync.WaitGroup
	codes := make([]int, concurrent)
	for i := 0; i < concurrent; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			req := httptest.NewRequest(http.MethodGet, "/api/incidents?concelho=Faro", nil)
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)
			codes[i] = rec.Code
		}(i)
	}

	close(client.release)
	wg.Wait()

	for i, code := range codes {
		if code != http.StatusOK {
			t.Errorf("request %d: want 200, got %d", i, code)
		}
	}
	if got := atomic.LoadInt32(&client.calls); got != 1 {
		t.Errorf("want 1 upstream call for %d concurrent requests, got %d", concurrent, got)
	}
}

func TestConcelhos(t *testing.T) {
	mux := handler.NewMux(newCountingClient())
	req := httptest.NewRequest(http.MethodGet, "/api/concelhos", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Lisboa") {
		t.Errorf("concelhos list missing expected entry: %s", rec.Body.String())
	}
	if cc := rec.Header().Get("Cache-Control"); cc != "public, max-age=86400" {
		t.Errorf("unexpected Cache-Control: %q", cc)
	}
}

func TestStaticHandler_ServesIndex(t *testing.T) {
	mux := handler.NewMux(newCountingClient())
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Fogos Ativos") {
		t.Errorf("index body missing expected title: %s", rec.Body.String())
	}
}

func TestSecurityHeaders(t *testing.T) {
	mux := handler.NewMux(newCountingClient())
	req := httptest.NewRequest(http.MethodGet, "/api/concelhos", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	want := map[string]string{
		"X-Content-Type-Options":    "nosniff",
		"Referrer-Policy":           "no-referrer",
		"X-Frame-Options":           "DENY",
		"Strict-Transport-Security": "max-age=31536000; includeSubDomains",
	}
	for header, expected := range want {
		if got := rec.Header().Get(header); got != expected {
			t.Errorf("header %s: want %q, got %q", header, expected, got)
		}
	}
	if rec.Header().Get("Content-Security-Policy") == "" {
		t.Error("missing Content-Security-Policy header")
	}
}

// errUpstream is a sentinel error simulating an upstream client failure.
var errUpstream = &testError{"upstream unavailable"}

type testError struct{ msg string }

func (e *testError) Error() string { return e.msg }
