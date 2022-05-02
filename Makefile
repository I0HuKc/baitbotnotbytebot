DOCKER_PATH=docker

.PHONY: run
run:
	docker-compose -f $(DOCKER_PATH)/docker-compose.yml build
	docker-compose -f $(DOCKER_PATH)/docker-compose.yml up


.PHONY: rund
rund:
	docker-compose -f $(DOCKER_PATH)/docker-compose.yml build
	docker-compose -f $(DOCKER_PATH)/docker-compose.yml up -d


.PHONY: build
build:
	docker-compose -f $(DOCKER_PATH)/docker-compose.yml build


.PHONY: down
down:
	docker-compose -f $(DOCKER_PATH)/docker-compose.yml down \
			--volumes \
			--remove-orphans


.PHONY: count
count:
	find . -name tests -prune -o -type f -name '*.go' | xargs wc -l

.DEFAULT_GOAL := run

