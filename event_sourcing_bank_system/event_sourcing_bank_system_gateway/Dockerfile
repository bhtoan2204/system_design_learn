FROM golang:1.24-alpine AS build

WORKDIR /usr/app

COPY event_sourcing_bank_system_gateway/go.mod .
COPY event_sourcing_bank_system_gateway/go.sum .
RUN go mod download

COPY event_sourcing_bank_system_gateway/ .

RUN go build -o /usr/app/event_sourcing_bank_system_gateway /usr/app/cmd/main.go

FROM alpine:3.19 AS production

WORKDIR /usr/app

COPY --from=build /usr/app ./

CMD ["./event_sourcing_bank_system_gateway"]