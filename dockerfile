FROM golang:1.26.3-alpine AS bulider
WORKDIR /app
COPY go.mod ./
RUN go.mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/hl-rate-limiter ./cmd/main.go
FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/hl-rate-limiter .
EXPOSE 8080
CMD ["./app/hl-rate-limiter"]