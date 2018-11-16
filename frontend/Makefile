build-local:
	grunt

build-prod:
	grunt prod

local:
	npm run local

docker-build-local: build-local
	docker build -t short-url-frontend .

docker-start-dev: docker-build-local
	docker run -d -p 80:80 short-url-frontend

docker-end-dev:
	docker stop short-url-frontend
	docker rm short-url-frontend
	dokcer rmi short-url-frontend
