# ghcr.io/okteto/golang has no Go 1.26.10 build yet. Switch back to ghcr.io/okteto/golang:1.26
# once `docker run --rm ghcr.io/okteto/golang:1.26 go version` prints go1.26.10 or later.
FROM golang:1.26.10 AS builder

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o bin/okteto-test -ldflags="-s -w" -trimpath ./cmd/main.go

FROM ghcr.io/okteto/okteto:3.24.1 AS final
WORKDIR /root/
COPY --from=builder /app/bin/okteto-test /okteto-test
ENTRYPOINT ["/okteto-test"]
