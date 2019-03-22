#Version: 0.0.1
FROM ubuntu:14.04
MAINTAINER calmwu "calm.wu@qq.com"
RUN apt-get update
RUN apt-get install -y nginx
RUN echo 'Hi, i am in your container' > /usr/share/nginx/html/index.html
EXPOSE  7878