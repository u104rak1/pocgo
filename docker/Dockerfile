# Build stage
FROM golang:1.23.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -o pocgo cmd/pocgo/main.go

# Runtime stage
FROM alpine:latest AS runtime

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/pocgo .

CMD ["./pocgo"]
