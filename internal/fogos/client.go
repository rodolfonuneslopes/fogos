package fogos

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const userAgent = "fogos-concelho-spa/1.0 (github.com/rodolfonuneslopes/fogos)"

// Incident represents a single active wildfire occurrence.
// Field names are confirmed against the fogos.pt v2 API response.
type Incident struct {
	ID             string `json:"id"`
	Concelho       string `json:"concelho"`
	Freguesia      string `json:"freguesia"`
	District       string `json:"district"`
	Status         string `json:"status"`
	Natureza       string `json:"natureza"`
	Date           string `json:"date"`
	Hour           string `json:"hour"`
	Man            int    `json:"man"`
	Terrain        int    `json:"terrain"`
	Aerial         int    `json:"aerial"`
	MeiosAquaticos int    `json:"meios_aquaticos"`
	StatusCode     int    `json:"statusCode"`
	DetailLocation string `json:"detailLocation"`
	Extra          string `json:"extra"`
}

// Client is the interface for fetching fogos.pt data.
type Client interface {
	ActiveIncidents(concelho string) ([]Incident, error)
}

type client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

// New returns a Client backed by the fogos.pt API.
func New(baseURL, token string) Client {
	return &client{
		baseURL: baseURL,
		token:   token,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *client) ActiveIncidents(concelho string) ([]Incident, error) {
	u, err := url.Parse(c.baseURL + "/v2/incidents/active")
	if err != nil {
		return nil, fmt.Errorf("parse base URL: %w", err)
	}
	q := u.Query()
	if concelho != "" {
		q.Set("concelho", concelho)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	if c.token != "" {
		req.Header.Set("FOGOS-PT-AUTH", c.token)
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("upstream request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("upstream returned %d", resp.StatusCode)
	}

	var result struct {
		Data []Incident `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	// nil → empty slice so callers can always range without a nil check
	if result.Data == nil {
		result.Data = []Incident{}
	}
	return result.Data, nil
}
