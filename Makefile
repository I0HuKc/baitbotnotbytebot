DOKCER_FILE_ENV=local
DOCKER_PATH=docker

.PHONY: run
run:
	docker-compose -f $(DOCKER_PATH)/docker-compose.$(DOKCER_FILE_ENV).yml build
	docker-compose -f $(DOCKER_PATH)/docker-compose.$(DOKCER_FILE_ENV).yml up

.PHONY: rund
rund:
	docker-compose -f $(DOCKER_PATH)/docker-compose.$(DOKCER_FILE_ENV).yml build
	docker-compose -f $(DOCKER_PATH)/docker-compose.$(DOKCER_FILE_ENV).yml up -d


.PHONY: build
build:
	docker-compose -f $(DOCKER_PATH)/docker-compose.$(DOKCER_FILE_ENV).yml build


.PHONY: down
down:
	docker-compose -f $(DOCKER_PATH)/docker-compose.$(DOKCER_FILE_ENV).yml down \
			--volumes \
			--remove-orphans

.DEFAULT_GOAL := run

