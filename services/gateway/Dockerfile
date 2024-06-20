FROM golang:1.22.4-alpine3.19 AS builder

RUN apk update --no-cache && apk add --no-cache tzdata

ENV CGO_ENABLED=0

WORKDIR /app/gateway

ADD go.mod .
ADD go.sum .

RUN go mod download

COPY . .

RUN go build -ldflags="-s -w" -o ./bin/gateway ./cmd/gateway/main.go

FROM alpine:3.19

COPY --from=builder /usr/share/zoneinfo/Europe/Moscow /usr/share/zoneinfo/Europe/Moscow
ENV TZ=Europe/Moscow

WORKDIR /app/bin

COPY --from=builder /app/gateway/bin .

ENTRYPOINT [ "./gateway" ]