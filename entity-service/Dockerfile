FROM golang:1.12.5-alpine3.9 AS build
RUN apk --no-cache add gcc g++ make ca-certificates
WORKDIR /go/src/github.com/lisin-andrey/msvc-demo

COPY common common
COPY entity-service entity-service
COPY vendor vendor
COPY ./entity-service/entity-service.yaml /go/bin

RUN go install ./...

FROM alpine:3.9
WORKDIR /usr/bin
COPY --from=build /go/bin . 