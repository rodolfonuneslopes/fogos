package fogos_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rodolfonuneslopes/fogos/internal/fogos"
)

func serve(t *testing.T, status int, body string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
}

func TestActiveIncidents_WithFires(t *testing.T) {
	ts := serve(t, 200, `{"data":[{"id":"1","concelho":"Castelo Branco","status":"Resolução"}]}`)
	defer ts.Close()

	c := fogos.New(ts.URL, "test-token")
	incidents, err := c.ActiveIncidents("Castelo Branco")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(incidents) != 1 {
		t.Fatalf("want 1 incident, got %d", len(incidents))
	}
	if incidents[0].Concelho != "Castelo Branco" {
		t.Errorf("wrong concelho: %q", incidents[0].Concelho)
	}
}

func TestActiveIncidents_NoFires(t *testing.T) {
	ts := serve(t, 200, `{"data":[]}`)
	defer ts.Close()

	c := fogos.New(ts.URL, "test-token")
	incidents, err := c.ActiveIncidents("Castelo Branco")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(incidents) != 0 {
		t.Fatalf("want 0 incidents, got %d", len(incidents))
	}
}

func TestActiveIncidents_NullData(t *testing.T) {
	ts := serve(t, 200, `{"data":null}`)
	defer ts.Close()

	c := fogos.New(ts.URL, "test-token")
	incidents, err := c.ActiveIncidents("Castelo Branco")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if incidents == nil {
		t.Fatal("expected non-nil slice for null data")
	}
}

func TestActiveIncidents_UpstreamError(t *testing.T) {
	ts := serve(t, 503, `Service Unavailable`)
	defer ts.Close()

	c := fogos.New(ts.URL, "test-token")
	_, err := c.ActiveIncidents("Castelo Branco")
	if err == nil {
		t.Fatal("expected error for non-200 status")
	}
}

func TestActiveIncidents_MalformedJSON(t *testing.T) {
	ts := serve(t, 200, `{bad json`)
	defer ts.Close()

	c := fogos.New(ts.URL, "test-token")
	_, err := c.ActiveIncidents("Castelo Branco")
	if err == nil {
		t.Fatal("expected error for malformed JSON")
	}
}

// captureHeaders serves an empty incident list and records the request headers.
func captureHeaders(t *testing.T, got *http.Header) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*got = r.Header.Clone()
		_, _ = w.Write([]byte(`{"data":[]}`))
	}))
}

// Both auth headers must be sent: X-API-KEY is what the upstream gateway
// enforces, FOGOS-PT-AUTH is what fogos.pt documents.
func TestActiveIncidents_SendsBothAuthHeaders(t *testing.T) {
	var got http.Header
	ts := captureHeaders(t, &got)
	defer ts.Close()

	c := fogos.New(ts.URL, "test-token")
	if _, err := c.ActiveIncidents("Castelo Branco"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, h := range []string{"X-API-KEY", "FOGOS-PT-AUTH"} {
		if v := got.Get(h); v != "test-token" {
			t.Errorf("%s = %q, want %q", h, v, "test-token")
		}
	}
}

func TestActiveIncidents_NoTokenSendsNoAuthHeaders(t *testing.T) {
	var got http.Header
	ts := captureHeaders(t, &got)
	defer ts.Close()

	c := fogos.New(ts.URL, "")
	if _, err := c.ActiveIncidents("Castelo Branco"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, h := range []string{"X-API-KEY", "FOGOS-PT-AUTH"} {
		if v := got.Get(h); v != "" {
			t.Errorf("%s = %q, want it unset", h, v)
		}
	}
}
