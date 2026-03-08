FROM golang:1.26.0-alpine3.23 AS build

WORKDIR /usr/src/app

COPY go.mod go.mod
COPY cmd cmd
COPY internal internal

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o /usr/local/bin/super-simple-queues ./cmd/super-simple-queues

FROM scratch

COPY --from=build /usr/local/bin/super-simple-queues /usr/local/bin/super-simple-queues

CMD ["/usr/local/bin/super-simple-queues"]