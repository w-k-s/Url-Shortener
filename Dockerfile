FROM alpine
MAINTAINER W.K.S <waqqas.abdulkareem@gmail.com>

RUN mkdir /app 
COPY build/ /app/ 
WORKDIR /app 
CMD ["/app/app"]

EXPOSE 8080