clean:
	rm -rf build/

fmt:
	gofmt -w .

run: fmt
	DB_CONN_STRING=postgresql://localhost:5432/url_shortener?sslmode=disable BASE_URL=http://localhost:8080 ADDRESS=localhost:8080 go run *.go

test: fmt
	go clean -testcache; go test -count=1 ./...

docker-build: fmt test clean
	docker build --no-cache -t wkas/short-url:dev .
