FROM golang:1.23.10-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o taskapi cmd/api/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/taskapi .
EXPOSE 8080
CMD ["./taskapi"]