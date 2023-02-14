# syntax=docker/dockerfile:1
# Step 1: build golang binary
FROM golang:1.19 as builder
WORKDIR /opt/app
ENV GOPROXY https://goproxy.cn,direct
COPY go.* ./
RUN go mod download
COPY . .
RUN go build -o example

# Step 2: copy binary from step1
FROM ubuntu:latest
WORKDIR /opt/app
COPY --from=builder /opt/app/example ./
CMD [ "/opt/app/example" ]
# 多阶段构建，同时使用Docker构建缓存