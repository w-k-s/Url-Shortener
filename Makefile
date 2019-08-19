clean:
	rm -rf build/

fmt:
	gofmt -w .

run: fmt
	DB_CONN_STRING=mongodb://localhost:27017/shorturl go run *.go

test: fmt
	go test ./...

docker-build: fmt test clean
	docker build --no-cache -t wkas/short-url:dev .

docker-start-dev: 
	docker-compose -f docker-config/docker-compose.dev.yml up -d

docker-stop-dev:
	docker-compose -f docker-config/docker-compose.dev.yml stop && docker-compose -f docker-config/docker-compose.dev.yml rm -f

docker-start-prod:
	docker-compose -f docker-config/docker-compose.production.yml up -d

docker-stop-prod:
	docker-compose -f docker-config/docker-compose.production.yml stop
	docker-compose -f docker-config/docker-compose.production.yml rm

docker-restart-prod: docker-stop-prod docker-start-prod

