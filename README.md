# fogos

An SPA to show active wildfires in a given Portuguese city (concelho), powered by [fogos.pt](https://fogos.pt) / ANEPC data.

## Running locally

The fogos.pt API is public and requires no token. The simplest way to run:

```bash
make run
```

The server listens on `:8888` by default (`LISTEN_ADDR` to override).

## Running with Docker

Pull the published image and run it directly — no Go toolchain needed:

```bash
docker run --rm -p 8888:8888 ghcr.io/rodolfonuneslopes/fogos:latest
```

Or build it from source:

```bash
docker build -t fogos:local .
docker run --rm -p 8888:8888 fogos:local
```

Pass any of the [Modes](#modes) environment variables with `-e`, e.g.
`-e FOGOS_MOCK=true`.

## Running with Docker Compose

```bash
docker compose up -d
```

See [docker-compose.yml](docker-compose.yml) — it pulls the published image
by default; uncomment `build: .` to build from source instead. Edit the
`environment:` block to switch modes or set a `FOGOS_TOKEN`.

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
