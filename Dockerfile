FROM golang:1.22 AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o okteto-cli ./cmd/main.go

FROM scratch
WORKDIR /root/
COPY --from=builder /app/okteto-cli .
CMD ["./okteto-cli"]
