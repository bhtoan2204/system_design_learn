from golang:1.24.1-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o main main.go

FROM alpine:latest

COPY --from=builder /app/main /app/main