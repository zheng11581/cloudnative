FROM golang:1.17.6 as build
WORKDIR /build
COPY . .
RUN go build -o httpserver .

FROM alpine:3.15
COPY --from=build /build/httpserver /httpserver
EXPOSE 80
ENTRYPOINT ["/httpserver"]