export tag=v1.0

root:
	export ROOT=cncamp

build:
	echo "building busyloop binary"
	mkdir -p bin/amd64
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/amd64

release: build
	echo "building busyloop container"
	docker build -tag zheng11581/busyloop:${tag}

push: release
	echo "push zheng11581/busyloop:${tag} to Dockerhub"