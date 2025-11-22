# Stage 1: Build Go binary
FROM golang:1.24 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o ngoclam-zmp-be ./cmd/main.go

# Stage 2: Final image (Alpine)
FROM alpine:3.20
RUN apk add --no-cache bash curl libc6-compat

WORKDIR /app

# Copy binary app
COPY --from=builder /app/ngoclam-zmp-be ./ngoclam-zmp-be

ENTRYPOINT ["./ngoclam-zmp-be"]
