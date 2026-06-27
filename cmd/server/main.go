package main

import (
	"log"
	"net/http"
	"os"

	"github.com/rodolfonuneslopes/fogos/internal/handler"
)

func main() {
	addr := envOrDefault("LISTEN_ADDR", ":8080")
	baseURL := envOrDefault("FOGOS_BASE_URL", "https://api.fogos.pt")
	token := os.Getenv("FOGOS_TOKEN")
	if token == "" {
		log.Fatal("FOGOS_TOKEN environment variable is required")
	}

	mux := handler.NewMux(baseURL, token)

	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
