BUILD_NAME = short-url

clean:
	rm -f $(BUILD_NAME)

fmt:
	gofmt -w .

run: fmt
	go run *.go

test: fmt
	go test ./...

docker-run-local: clean fmt
	go build 
	docker-compose -f docker-compose.local.yml build
	docker-compose -f docker-compose.local.yml up -d

docker-end-local:
	docker-compose -f docker-compose.local.yml stop
	docker-compose -f docker-compose.local.yml rm

docker-hub-publish: clean fmt test
	docker-compose -f docker-compose.production.yml build
	docker-compose -f docker-compose.production.yml push short-url