BUILD_NAME = short-url

clean:
	rm -f $(BUILD_NAME)

fmt:
	gofmt -w .

dep:
	godep save

run: dep fmt
	go run *.go

test: fmt
	#Ignore the vendor directory
	go test  `go list ./... | grep -v vendor`

docker-run-local: clean dep fmt
	go build 
	docker-compose -f docker-compose.local.yml build
	docker-compose -f docker-compose.local.yml up -d

docker-end-local:
	docker-compose -f docker-compose.local.yml stop
	docker-compose -f docker-compose.local.yml rm

docker-hub-publish: clean dep fmt test
	docker-compose -f docker-compose.production.yml build
	docker-compose -f docker-compose.production.yml push short-url