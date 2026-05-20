# Build stage
FROM golang:alpine AS build-env

ADD . /go/src/github.com/pruh/api/v3/

RUN apk update \
    && apk add --no-cache \
    gcc \
    musl-dev

RUN cd /go/src/github.com/pruh/api/v3/ && go build -o api

# Run stage
FROM alpine
RUN apk update \
    && apk add ca-certificates
WORKDIR /app
COPY --from=build-env /go/src/github.com/pruh/api/v3/api /app/

CMD ["/app/api"]
