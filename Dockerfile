FROM --platform=$BUILDPLATFORM golang:1.27-alpine@sha256:b8c9b82affbbb82a020bbcf998dad2ba9516944ac7e883524642a481a9505e9c AS build

ARG TARGETARCH
ARG VERSION=dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH go build \
    -ldflags "-s -w -X yaddc/internal/yaddc.Version=${VERSION}" \
    -o yaddc ./main.go

FROM gcr.io/distroless/static-debian12@sha256:d75cdd72874d4790092fcb1b058493ecf6bb5bf2b2b897045b00ff01d91843f2

COPY --from=build /app/yaddc /yaddc

USER nonroot:nonroot

LABEL org.opencontainers.image.source="https://github.com/daanvzijl/yet-another-dynamic-dns-client"
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.description="A lightweight dynamic DNS client"

ENTRYPOINT ["/yaddc"]
