FROM golang:1.16.6-buster AS builder

ENV GO111MODULE=on
ENV GOOS=linux
ENV GOARCH=amd64
WORKDIR /go/src/github.com/dai65527/microservice-handson

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . .
RUN go build -o /usr/bin/db ./platform/db

## Runtime

FROM gcr.io/distroless/base:1e4a8bb3ad03f71b572cbcb3bbc25f3fd8d0ff14
COPY --from=builder /usr/bin/db /usr/bin/db

CMD ["/usr/bin/db"]
