ENV=local
DOCKER_PATH=docker

.PHONY: run
run:
	docker-compose -f $(DOCKER_PATH)/docker-compose.$(ENV).yml build
	docker-compose -f $(DOCKER_PATH)/docker-compose.$(ENV).yml up


.PHONY: build
build:
	docker-compose -f $(DOCKER_PATH)/docker-compose.$(ENV).yml build


.PHONY: down
down:
	docker-compose -f $(DOCKER_PATH)/docker-compose.$(ENV).yml down \
			--volumes \
			--remove-orphans

.DEFAULT_GOAL := run

