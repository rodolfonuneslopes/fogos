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
	return mux
}

func incidentsHandler(client fogos.Client, cache *incidentCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		concelho := r.URL.Query().Get("concelho")

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
