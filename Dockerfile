FROM golang:1.20-alpine AS builder
# Dependencies for go-sqlite3
ENV CGO_ENABLED=1
RUN apk add --no-cache gcc musl-dev

WORKDIR /app
COPY . .
RUN go build -o main

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main /usr/local/bin/monolith
ENTRYPOINT ["/usr/local/bin/monolith", "-serve"]

