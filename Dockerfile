# Stage 1: Build the Svelte frontend
FROM node:22-alpine AS frontend
WORKDIR /frontend
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

# Stage 2: Build the Go binary
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend /web/static ./web/static
RUN CGO_ENABLED=0 GOOS=linux go build -o harmonic .

# Stage 3: Final minimal image
FROM alpine:3.21
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/harmonic .
COPY --from=builder /app/web/static ./web/static

EXPOSE 8080 22222
CMD ["./harmonic", "--port", "22222", "--http", "8080"]
