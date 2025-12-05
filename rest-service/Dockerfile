FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o notes-api ./cmd

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/notes-api /app/notes-api
COPY .env.dev /app/.env.dev

EXPOSE 8080 50051

CMD ["/app/notes-api"]
