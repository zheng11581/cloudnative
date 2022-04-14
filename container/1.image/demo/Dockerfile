FROM ubuntu:impish

RUN apt update && apt-get install -y nginx

COPY file1 /var/www/html/
ADD  file2.tar.gz /var/www/html/
CMD ["/usr/sbin/nginx", "-g", "daemon off;"]
