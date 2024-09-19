FROM golang:1.22 AS builder

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o bin/okteto-test -ldflags="-s -w" ./cmd/main.go

FROM okteto/okteto:2.31.0 AS final
WORKDIR /root/
COPY --from=builder /app/bin/okteto-test /okteto-test
ENTRYPOINT ["/okteto-test"]
