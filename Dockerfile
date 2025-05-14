FROM golang:1.23 as builder

WORKDIR /app

COPY go.mod go.sum ./

COPY . .
RUN go mod download

RUN go build -o user-service ./cmd/user-service

FROM debian:bookworm-slim

COPY --from=builder /app/user-service /usr/local/bin/user-service

CMD ["/usr/local/bin/user-service"]