# Yet another dynamic dns client

A lightweight dynamic DNS client that keeps your DNS A records in sync with your current public IP.

## Supported DNS providers

| Provider   | Environment variables        |
| ---------- | ---------------------------- |
| Cloudflare | `CF_API_TOKEN`, `CF_ZONE_ID` |

## Setup

**Prerequisites**

- mise

Run `just install` to get going.

## Configuration

| Variable             | Required   | Description                                                                                  |
| -------------------- | ---------- | -------------------------------------------------------------------------------------------- |
| `YADDC_A_RECORDS`    | Yes        | Comma-separated list of A records to sync, e.g. `record-a.example.com, record-b.example.com` |
| `YADDC_DNS_PROVIDER` | Yes        | DNS backend to use. Defaults to `cloudflare`                                                 |
| `CF_API_TOKEN`       | Cloudflare | Cloudflare API token with DNS edit permissions                                               |
| `CF_ZONE_ID`         | Cloudflare | Cloudflare zone ID for your domain                                                           |

## Run

```sh
YADDC_DNS_PROVIDER=cloudflare \
CF_API_TOKEN=your_token \
CF_ZONE_ID=your_zone_id \
YADDC_A_RECORDS=record-a.example.com \
go run .
```

## Docker

```sh
docker run \
  -e YADDC_DNS_PROVIDER=cloudflare \
  -e CF_API_TOKEN=your_token \
  -e CF_ZONE_ID=your_zone_id \
  -e YADDC_A_RECORDS=record-a.example.com \
  ghcr.io/daanvzijl/yet-another-dynamic-dns-client:latest
```

Multi-platform images are published for `linux/amd64` and `linux/arm64`.

## Pre-built binaries

Binaries for Linux, macOS, and Windows are attached to each [release](../../releases).

| File                      | Platform            |
| ------------------------- | ------------------- |
| `yaddc-linux-amd64`       | Linux x86_64        |
| `yaddc-linux-arm64`       | Linux arm64         |
| `yaddc-darwin-amd64`      | macOS x86_64        |
| `yaddc-darwin-arm64`      | macOS Apple Silicon |
| `yaddc-windows-amd64.exe` | Windows x86_64      |

## Development

```sh
just test       # run tests
just test-race  # run tests with race detector
just lint       # run all linters
```

## Release

Releases are triggered by merging a PR with one of the following labels:

| Label           | Effect                              |
| --------------- | ----------------------------------- |
| `release:patch` | Bumps patch version (1.0.0 → 1.0.1) |
| `release:minor` | Bumps minor version (1.0.0 → 1.1.0) |
| `release:major` | Bumps major version (1.0.0 → 2.0.0) |

If multiple labels are applied, `major` takes precedence over `minor` over `patch`.
