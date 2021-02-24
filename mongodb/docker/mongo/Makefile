all: clean build

.PHONY: build
build:
	docker-compose up -d

.PHONY: stop
stop:
	docker-compose down --remove-orphans -v

.PHONY: primary
primary:
	docker-compose exec mongo1 mongo --port 9142

.PHONY: secondary
secondary:
	docker-compose exec mongo2 mongo --port 9242

.PHONY: clean
clean: stop
	rm -rf data*