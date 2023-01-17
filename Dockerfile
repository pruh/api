# Build stage
FROM golang:1.19-alpine AS build-env

ADD . /go/src/github.com/pruh/api/

RUN apk update \
    && apk add --no-cache \
    gcc \
    musl-dev

RUN cd /go/src/github.com/pruh/api/ && go build -o api

# Run stage
FROM alpine
RUN apk update \
    && apk add ca-certificates
WORKDIR /app
COPY --from=build-env /go/src/github.com/pruh/api/api /app/

CMD ["/app/api"]
