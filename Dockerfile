# syntax=docker/dockerfile:1
FROM golang:1.19 as build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY pkg pkg
COPY cmd cmd

RUN go build -ldflags="-w -s" -o /jsplit ./cmd/jsplit

FROM debian:bullseye-slim as runtime
COPY --from=build /jsplit /jsplit
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

CMD [ "/jsplit" ]