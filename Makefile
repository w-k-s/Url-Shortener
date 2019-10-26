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

