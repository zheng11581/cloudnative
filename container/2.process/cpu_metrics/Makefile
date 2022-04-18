clean:
	rm -f *.o threads-cpu
	docker rmi registry/cpu_metrics:v1
threads-cpu: clean
	gcc -o threads-cpu threads-cpu.c -lpthread
image: threads-cpu
	docker build -t registry/cpu_metrics:v1 .
all: image
