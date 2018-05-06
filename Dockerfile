# Build stage
FROM golang:1.9.3-alpine AS build-env

ADD . /go/src/github.com/pruh/api/

RUN apk update \
    && apk add --no-cache git \
    && go get github.com/gorilla/mux \
    && go get github.com/urfave/negroni \
    && apk del git

RUN cd /go/src/github.com/pruh/api/ && go build -o api

# Run stage
FROM alpine
RUN apk update \
    && apk add ca-certificates
WORKDIR /app
COPY --from=build-env /go/src/github.com/pruh/api/api /app/
