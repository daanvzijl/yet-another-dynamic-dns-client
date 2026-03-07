# yaddc

Yet another dynamic DNS client. Syncs Cloudflare A records to your current public IP.

## Prerequisites

- [mise](https://mise.jdx.dev/)

## Setup

```sh
just install
```

## Configuration

| Variable          | Required | Description                                                                                  |
| ----------------- | -------- | -------------------------------------------------------------------------------------------- |
| `CF_API_TOKEN`    | Yes      | Cloudflare API token with DNS edit permissions                                               |
| `CF_ZONE_ID`      | Yes      | Cloudflare zone ID                                                                           |
| `YADDC_A_RECORDS` | Yes      | Comma-separated list of A records to sync (e.g. `record-a.example.com,record-b.example.com`) |

## Usage

```sh
just run
```

Or build and run:

```sh
just build
```

## Docker

```sh
just docker-build
docker run --rm \
  -e CF_API_TOKEN=... \
  -e CF_ZONE_ID=... \
  -e YADDC_A_RECORDS=home.example.com \
  yaddc
```

## Development

```sh
just lint       # run all linters
just fix        # run linters with auto-fix
just test-race  # run tests with race detector
```

## License

MIT
