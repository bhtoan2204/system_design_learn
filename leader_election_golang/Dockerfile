FROM golang:1.24.1

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build -o main ./src

CMD ["./main"]