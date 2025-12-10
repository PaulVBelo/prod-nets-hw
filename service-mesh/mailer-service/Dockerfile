FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o mailer ./cmd/main.go

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/mailer /app/mailer

EXPOSE 8080

CMD ["/app/mailer"]
