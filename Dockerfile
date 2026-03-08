FROM golang:1.26-alpine AS build

ARG TARGETARCH
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH go build -o yaddc ./main.go

FROM gcr.io/distroless/static-debian12

COPY --from=build /app/yaddc /yaddc

LABEL org.opencontainers.image.source="https://github.com/daanvzijl/yet-another-dynamic-dns-client"
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.description="A lightweight dynamic DNS client"

ENTRYPOINT ["/yaddc"]
