fmt:
	gofmt -w .

run: fmt
	go run *.go

run-local-docker: fmt
	go build 
	docker-compose build
	docker-compose up -d

end-local-docker:
	docker-compose stop
	docker-compose rm