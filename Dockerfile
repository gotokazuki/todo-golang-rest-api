# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o todo-app

FROM gcr.io/distroless/static-debian12:latest
WORKDIR /app
COPY --from=builder /app/todo-app .
EXPOSE 8080
CMD ["./todo-app"]
