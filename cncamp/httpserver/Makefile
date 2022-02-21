export tag=v1.0
root:
	export ROOT=cncamp

build:
	echo "building httpserver binary"
	#CGO_ENABLED=0 GOOS=maos GOARCH=amd64 go build -o .
	go build -o .

release: build
	echo "building httpserver container"
	docker build -t zheng11581/httpserver:${tag} .

push: release
	echo "pushing cncamp/httpserver"
	docker push zheng11581/httpserver:v1.0