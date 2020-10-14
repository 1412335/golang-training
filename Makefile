.PHONY: run
run:
	docker build -t golang-traning:week2 .
	docker run --rm --name golang-traning-week2 -p 8080:8080 golang-traning:week2

.PHONY: k6-test
k6-test:
	docker run -i loadimpact/k6 run --vus 10 --duration 30s - <week2/script.js

.PHONY: week2
week2:
	cd ./docker/k6-docker && \
	docker-compose up -d influxdb grafana goapp

.PHONY: test
test:
	cd ./docker/k6-docker && \
	docker-compose run -v `pwd`/../../week2:/scripts k6 run --vus 10 --duration 30s /scripts/script.js