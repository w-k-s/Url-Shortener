FROM golang:1.12.4 as builder

WORKDIR /go/src/github.com/w-k-s/short-url

COPY . /go/src/github.com/w-k-s/short-url

ENV GO111MODULE=on

RUN go get

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o short-url *.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /go/src/github.com/w-k-s/short-url/short-url .

CMD ["./short-url"]