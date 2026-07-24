FROM golang:1.26-alpine@sha256:0178a641fbb4858c5f1b48e34bdaabe0350a330a1b1149aabd498d0699ff5fb2 AS build

ARG TARGETARCH
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH go build -o yaddc ./main.go

FROM gcr.io/distroless/static-debian12@sha256:a9fcaedd4c9b59e12dd65d954f0b5044f19b0647a8a3712e77205df9e7b102cd

COPY --from=build /app/yaddc /yaddc

LABEL org.opencontainers.image.source="https://github.com/daanvzijl/yet-another-dynamic-dns-client"
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.description="A lightweight dynamic DNS client"

ENTRYPOINT ["/yaddc"]
