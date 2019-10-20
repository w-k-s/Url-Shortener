clean:
	rm -rf build/

fmt:
	gofmt -w .

run: fmt
	DB_CONN_STRING=postgresql://localhost:5432/url_shortener?sslmode=disable go run *.go

test: fmt
	go clean -testcache; go test -count=1 ./...

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

