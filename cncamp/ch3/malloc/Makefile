root:
	mkdir -p bin/linux
build: root
	CGO_ENABLED=1 GOOS=linux CC=x86_64-linux-musl-gcc  CXX=x86_64-linux-musl-g++ go build -o bin/linux