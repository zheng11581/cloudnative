FROM golang:1.17.6 as build
ENV GO111MODULE=off \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
WORKDIR /build
COPY . .
RUN go build -o httpserver

FROM alpine:3.15
COPY --from=build /build/httpserver /httpserver
EXPOSE 8080

RUN apk add --no-cache tini
# Tini is now available at /sbin/tini
ENTRYPOINT ["/sbin/tini", "--"]

# Run your program under Tini
 CMD ["/httpserver"]
# ENTRYPOINT ["/httpserver"]