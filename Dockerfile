FROM golang:1.9.2
MAINTAINER W.K.S <waqqas.abdulkareem@gmail.com>

RUN mkdir -p /go/src/github.com/w-k-s/short-url
WORKDIR /go/src/github.com/w-k-s/short-url

COPY . /go/src/github.com/w-k-s/short-url
RUN go-wrapper download && go-wrapper install

CMD ["go-wrapper", "run"]

EXPOSE 8080