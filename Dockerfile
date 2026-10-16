FROM golang:1.27-alpine@sha256:f92b6ef800e499660581efdabdf25d9d817a9d124eaf900924f0504e7e27e12d AS build

ARG TARGETARCH
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH go build -o yaddc ./main.go

FROM gcr.io/distroless/static-debian12@sha256:d75cdd72874d4790092fcb1b058493ecf6bb5bf2b2b897045b00ff01d91843f2

COPY --from=build /app/yaddc /yaddc

LABEL org.opencontainers.image.source="https://github.com/daanvzijl/yet-another-dynamic-dns-client"
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.description="A lightweight dynamic DNS client"

ENTRYPOINT ["/yaddc"]
