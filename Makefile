clean:
	rm -rf build/

fmt:
	gofmt -w .

dep:
	godep save

run: fmt dep 
	go run *.go

test: fmt
	#Ignore the vendor directory
	go test  `go list ./... | grep -v vendor`

docker-build: fmt dep test clean
	docker build --no-cache -t wkas/short-url:dev .

docker-start-dev: docker-build
	docker-compose -f docker-compose.dev.yml up -d

docker-stop-dev:
	docker-compose -f docker-compose.dev.yml stop
	docker-compose -f docker-compose.dev.yml rm

docker-start-prod:
	docker-compose -f docker-compose.production.yml up -d

docker-stop-prod:
	docker-compose -f docker-compose.production.yml stop
	docker-compose -f docker-compose.production.yml rm

docker-restart-prod: docker-stop-prod docker-start-prod

