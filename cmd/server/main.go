package main

import (
	"log"
	"net/http"
	"os"

	"github.com/rodolfonuneslopes/fogos/internal/fogos"
	"github.com/rodolfonuneslopes/fogos/internal/handler"
)

func main() {
	addr := envOrDefault("LISTEN_ADDR", ":8080")

	var client fogos.Client
	if os.Getenv("FOGOS_MOCK") == "true" {
		log.Println("running in mock mode — no real API calls will be made")
		client = fogos.MockClient{}
	} else {
		token := os.Getenv("FOGOS_TOKEN")
		if token == "" {
			log.Fatal("FOGOS_TOKEN environment variable is required (or set FOGOS_MOCK=true)")
		}
		baseURL := envOrDefault("FOGOS_BASE_URL", "https://api.fogos.pt")
		client = fogos.New(baseURL, token)
	}

	mux := handler.NewMux(client)

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
