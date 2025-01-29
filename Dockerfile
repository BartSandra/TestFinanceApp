FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git build-base
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o /app/main ./cmd

FROM alpine:latest

RUN apk add --no-cache libpq bash curl

COPY --from=builder /app /app
COPY db/migrations /app/migrations
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY wait-for-it.sh /usr/local/bin/wait-for-it

RUN chmod +x /usr/local/bin/wait-for-it

WORKDIR /app

EXPOSE 8080

CMD ["sh", "-c", "/usr/local/bin/wait-for-it db:5432 --strict --timeout=30 -- goose -dir /app/migrations postgres \"postgres://postgres:postgres@db:5432/testdb?sslmode=disable\" up && /app/main"]


