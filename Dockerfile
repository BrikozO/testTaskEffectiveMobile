FROM golang:1.25 AS builder

RUN go install github.com/swaggo/swag/cmd/swag@latest

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN swag init -g main.go -o ./docs
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080