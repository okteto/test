FROM golang:1.23 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o okteto-action main.go

FROM okteto/okteto:2.31.0

COPY --from=builder /app/okteto-action .
RUN chmod +x ./okteto-action

ENTRYPOINT ["/entrypoint.sh"] 