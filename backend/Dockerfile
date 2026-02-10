# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o api ./cmd/api

# Runtime stage
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/api /api
COPY --from=builder /app/resume.pdf /resume.pdf
EXPOSE 8080
ENTRYPOINT ["/api"]
