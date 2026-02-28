FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /api ./cmd/api
RUN CGO_ENABLED=0 go build -o /seed ./cmd/seed-missions

FROM alpine:3.21
RUN apk add --no-cache ca-certificates
COPY --from=builder /api /usr/local/bin/api
COPY --from=builder /seed /usr/local/bin/seed
EXPOSE 8080
CMD ["api"]
