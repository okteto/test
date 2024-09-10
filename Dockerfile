FROM golang:1.23 as builder

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o okteto-action main.go
RUN chmod +x okteto-action

FROM okteto/okteto:2.31.0

COPY --from=builder /app/okteto-action .

ENTRYPOINT ["./okteto-action"] 