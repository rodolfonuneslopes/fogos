# fogos

An SPA to show active wildfires in a given Portuguese city (concelho), powered by [fogos.pt](https://fogos.pt) / ANEPC data.

## Running locally

The fogos.pt API is public and requires no token. The simplest way to run:

```bash
make run
```

The server listens on `:8080` by default (`LISTEN_ADDR` to override).

## Modes

Three modes are available, controlled entirely by environment variables:

| Mode | How to activate |
|---|---|
| Real API, no auth | (nothing set — the default) |
| Real API, with auth token | `FOGOS_TOKEN=your_token` |
| Mock (no network calls) | `FOGOS_MOCK=true` |

**Mock mode** returns fake incidents for every concelho and is useful for UI development without any network dependency. Most concelhos return two sample incidents; selecting **Évora** returns an empty list so you can see the "no active fires" state.

See [.env.example](.env.example) for a reference configuration.

## Available make targets

| Target | Description |
|--------|-------------|
| `make run` | Run the server |
| `make build` | Compile the binary |
| `make test` | Run all tests |
| `make fmt` | Format source with `gofmt` |
| `make vet` | Run `go vet` |
