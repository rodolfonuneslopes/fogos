package handler

import (
	"encoding/json"
	"net/http"

	"github.com/rodolfonuneslopes/fogos/internal/fogos"
)

// NewMux wires all routes and returns the root handler.
func NewMux(baseURL, token string) http.Handler {
	client := fogos.New(baseURL, token)
	cache := newIncidentCache()

	mux := http.NewServeMux()
	mux.Handle("GET /api/incidents", incidentsHandler(client, cache))
	mux.Handle("GET /api/concelhos", concelhosHandler())
	mux.Handle("/", http.FileServer(http.Dir("web")))
	return mux
}

func incidentsHandler(client fogos.Client, cache *incidentCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		concelho := r.URL.Query().Get("concelho")
		if concelho == "" {
			http.Error(w, "missing concelho parameter", http.StatusBadRequest)
			return
		}

		if cached, ok := cache.get(concelho); ok {
			writeJSON(w, cached)
			return
		}

		incidents, err := client.ActiveIncidents(concelho)
		if err != nil {
			http.Error(w, "upstream error", http.StatusBadGateway)
			return
		}
		cache.set(concelho, incidents)
		writeJSON(w, incidents)
	}
}

func concelhosHandler() http.HandlerFunc {
	data, _ := json.Marshal(concelhosList)
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, "encode error", http.StatusInternalServerError)
	}
}
