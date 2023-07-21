FROM golang:1.20 AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o outline-cli ./cli/main.go

FROM scratch

WORKDIR /app

COPY --from=builder /app/outline-cli .

ENTRYPOINT ["./outline-cli"]
