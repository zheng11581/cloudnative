FROM golang:1.17.6 as build
ENV GO111MODULE=off \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
WORKDIR /build
COPY . .
RUN make build

FROM alpine:3.15
COPY --from=build /build/httpserver /httpserver
EXPOSE 80
ENTRYPOINT ["/httpserver"]