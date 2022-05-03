DOCKER_DIR=docker

.PHONY: build
build:
	docker-compose -f $(DOCKER_DIR)/docker-compose.yml build


.PHONY: run
run: build	
	docker-compose -f $(DOCKER_DIR)/docker-compose.yml up


.PHONY: rund
rund: build	
	docker-compose -f $(DOCKER_DIR)/docker-compose.yml up -d


.PHONY: stop
stop:
	sudo docker stop \
			baitbot \
			baitbot_redis \
			baitbot_memcached


.PHONY: down
down:
	docker-compose -f $(DOCKER_DIR)/docker-compose.yml down \
			--volumes \
			--remove-orphans


.PHONY: count
count:
	find . -name tests -prune -o -type f -name '*.go' | xargs wc -l

.DEFAULT_GOAL := run

