include .env
export $(shell sed 's/=.*//' .env)

BUILD_NAME = short-url

clean:
	rm -f $(BUILD_NAME)

fmt:
	gofmt -w .

dep:
	godep save

run: fmt dep 
	go run *.go

test: fmt
	#Ignore the vendor directory
	go test  `go list ./... | grep -v vendor`

docker-build: fmt clean dep test
	docker build -t $(BUILD_NAME):$(TAG) 

docker-start-local: fmt clean dep 
	go build 
	docker-compose -f docker-compose.local.yml build
	docker-compose -f docker-compose.local.yml up -d

docker-end-local:
	docker-compose -f docker-compose.local.yml stop
	docker-compose -f docker-compose.local.yml rm

docker-start-prod:
	docker-compose -f docker-compose.production.yml up -d

docker-hub-publish: docker-build
	docker-compose -f docker-compose.production.yml push short-url