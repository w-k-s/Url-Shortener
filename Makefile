include .env
export $(shell sed 's/=.*//' .env)

BUILD_NAME = wkas/short-url

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
	docker build -t $(BUILD_NAME):$(TAG) .

docker-start-dev: fmt clean dep 
	go build 
	docker-compose -f docker-compose.dev.yml build
	docker-compose -f docker-compose.dev.yml up -d

docker-end-dev:
	docker-compose -f docker-compose.dev.yml stop
	docker-compose -f docker-compose.dev.yml rm

docker-start-prod:
	docker-compose -f docker-compose.production.yml up -d

docker-hub-publish: docker-build
	docker push $(BUILD_NAME)