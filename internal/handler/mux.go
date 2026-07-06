package handler

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"strconv"

	"github.com/rodolfonuneslopes/fogos/internal/fogos"
)

// NewMux wires all routes and returns the root handler.
func NewMux(client fogos.Client) http.Handler {
	cache := newIncidentCache()

	sub, _ := fs.Sub(webFS, "web")

	mux := http.NewServeMux()
	mux.Handle("GET /api/incidents", incidentsHandler(client, cache))
	mux.Handle("GET /api/concelhos", concelhosHandler())
	mux.Handle("/", staticHandler(sub))
	return securityHeaders(mux)
}

// securityHeaders sets standard hardening headers on every response. Public
// TLS is terminated at the Cloudflare Tunnel edge, not in this binary, but
// Strict-Transport-Security still reaches the browser through the tunnel.
func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("Referrer-Policy", "no-referrer")
		h.Set("Permissions-Policy", "geolocation=(), camera=(), microphone=()")
		h.Set("X-Frame-Options", "DENY")
		// img-src allows data: on top of 'self' because Pico CSS renders its
		// built-in icons (accordion chevron, checkbox, spinner, ...) as
		// inline data: URI SVG background-images; default-src alone blocks them.
		h.Set("Content-Security-Policy", "default-src 'self'; img-src 'self' data:; base-uri 'none'; frame-ancestors 'none'")
		h.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		next.ServeHTTP(w, r)
	})
}

func incidentsHandler(client fogos.Client, cache *incidentCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		concelho := r.URL.Query().Get("concelho")
		if concelho != "" && !validConcelhos[concelho] {
			http.Error(w, "unknown concelho", http.StatusBadRequest)
			return
		}

		if cached, ok := cache.get(concelho); ok {
			writeJSONBytes(w, cached)
			return
		}

		// singleflight: collapse concurrent cold-cache requests for the same
		// concelho into one upstream call; all waiters share the result.
		v, err, _ := cache.group.Do(concelho, func() (any, error) {
			incidents, err := client.ActiveIncidents(concelho)
			if err != nil {
				return nil, err
			}
			b, err := json.Marshal(incidents)
			if err != nil {
				return nil, fmt.Errorf("marshal incidents: %w", err)
			}
			cache.set(concelho, b)
			return b, nil
		})
		if err != nil {
			http.Error(w, "upstream error", http.StatusBadGateway)
			return
		}

		writeJSONBytes(w, v.([]byte))
	}
}

func concelhosHandler() http.HandlerFunc {
	data, _ := json.Marshal(concelhosList)
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "public, max-age=86400")
		w.Header().Set("Content-Length", strconv.Itoa(len(data)))
		w.Write(data)
	}
}

func staticHandler(fsys fs.FS) http.Handler {
	fileServer := http.FileServer(http.FS(fsys))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		fileServer.ServeHTTP(w, r)
	})
}

func writeJSONBytes(w http.ResponseWriter, b []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=60")
	w.Header().Set("Content-Length", strconv.Itoa(len(b)))
	w.Write(b)
}
