# fogos

An SPA to show active wildfires in a given Portuguese município (concelho), powered by [fogos.pt](https://fogos.pt) / ANEPC data.

## Running locally

Copy the example env file and fill in your fogos.pt API token:

```bash
cp .env.example .env
# edit .env and set FOGOS_TOKEN
```

Then:

```bash
make run
```

The server listens on `:8080` by default (`LISTEN_ADDR` to override).

## Mock mode

No API token yet? Run with fake data to develop and preview the UI:

```bash
FOGOS_MOCK=true make run
```

In mock mode no real API calls are made. Most concelhos return two sample incidents; selecting **Évora** returns an empty list so you can see the "no active fires" state.

## Available make targets

| Target | Description |
|--------|-------------|
| `make run` | Run the server (requires `.env`) |
| `make build` | Compile the binary |
| `make test` | Run all tests |
| `make fmt` | Format source with `gofmt` |
| `make vet` | Run `go vet` |
