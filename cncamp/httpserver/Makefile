export tag=v3.0-metrics
root:
	export ROOT=cncamp

build:
	echo "building httpserver binary"
	#CGO_ENABLED=0 GOOS=maos GOARCH=amd64 go build -o .
	go build -o .

release:
	echo "building httpserver container"
	docker build -t zheng11581/httpserver:${tag} .

push: release
	echo "pushing zheng11581/httpserver"
	docker push zheng11581/httpserver:${tag}