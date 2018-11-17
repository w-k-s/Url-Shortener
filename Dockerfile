FROM golang:1.9.2
MAINTAINER W.K.S <waqqas.abdulkareem@gmail.com>

RUN mkdir /app 
COPY build/ /app/ 
WORKDIR /app 
CMD ["/app/app"]

EXPOSE 8080