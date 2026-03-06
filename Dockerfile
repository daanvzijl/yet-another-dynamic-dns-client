FROM golang:1.25 as build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o yaddc ./main.go

FROM gcr.io/distroless/static-debian12

COPY --from=build /app/yaddc /yaddc
ENTRYPOINT ["/yaddc"]
