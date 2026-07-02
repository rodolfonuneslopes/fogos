package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/rodolfonuneslopes/fogos/internal/fogos"
	"github.com/rodolfonuneslopes/fogos/internal/handler"
)

func main() {
	addr := envOrDefault("LISTEN_ADDR", ":8888")

	var client fogos.Client
	if os.Getenv("FOGOS_MOCK") == "true" {
		log.Println("running in mock mode — no real API calls will be made")
		client = fogos.MockClient{}
	} else {
		token := os.Getenv("FOGOS_TOKEN")
		baseURL := envOrDefault("FOGOS_BASE_URL", "https://api.fogos.pt")
		if token != "" {
			log.Println("using real API with auth token")
		} else {
			log.Println("using real API without auth token")
		}
		client = fogos.New(baseURL, token)
	}

	srv := &http.Server{
		Addr:    addr,
		Handler: handler.NewMux(client),
		// ReadTimeout caps the time to read the full request including body.
		// 5s is generous for our header-only GET requests; it stops slow-loris attacks.
		ReadTimeout: 5 * time.Second,
		// WriteTimeout covers the time from end of request headers to end of response.
		// Matches the upstream client timeout (10s) plus a small buffer.
		WriteTimeout: 10 * time.Second,
		// IdleTimeout limits how long an idle keep-alive connection is held open.
		IdleTimeout: 120 * time.Second,
	}

	log.Printf("listening on %s", addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
