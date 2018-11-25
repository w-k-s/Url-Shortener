SHORT_URL_FRONTEND_DOCKER_ID=wkas
SHORT_URL_FRONTEND_IMAGE_NAME=short-url-frontend
SHORT_URL_FRONTEND_VERSION=1.0.1

BUILD_NAME=$(SHORT_URL_FRONTEND_DOCKER_ID)/$(SHORT_URL_FRONTEND_IMAGE_NAME)
IMAGE_NAME=$(BUILD_NAME):$(SHORT_URL_FRONTEND_VERSION)

clean:
	rm -rf build

local: clean
	npm run start

docker-build-dev: clean
	grunt
	docker build -t $(IMAGE_NAME) .

docker-start-dev: clean
	docker run --name $(SHORT_URL_FRONTEND_IMAGE_NAME) -d -p 8081:8081 $(IMAGE_NAME)

docker-build-prod: clean
	grunt production
	docker build -t $(IMAGE_NAME) .

docker-start-prod: docker-build-prod
	docker run --name $(SHORT_URL_FRONTEND_IMAGE_NAME) -d -p 8081:8081 $(IMAGE_NAME)

dockerhub-publish: docker-build-prod
	docker push $(IMAGE_NAME)

docker-end:
	docker stop $(SHORT_URL_FRONTEND_IMAGE_NAME)
	docker rm $(SHORT_URL_FRONTEND_IMAGE_NAME)
	dokcer rmi $(IMAGE_NAME)