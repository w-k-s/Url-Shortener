FROM golang:1.9.2 as builder
MAINTAINER W.K.S <waqqas.abdulkareem@gmail.com>

WORKDIR /go/src/github.com/w-k-s/short-url

COPY . /go/src/github.com/w-k-s/short-url
RUN go get -v -t -d ./...

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o short-url *.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /go/src/github.com/w-k-s/short-url/short-url .

CMD ["./short-url"]