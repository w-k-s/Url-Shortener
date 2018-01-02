fmt:
	gofmt -w .

run: fmt
	go run *.go

build: fmt
	go build

