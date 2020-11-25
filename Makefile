.PHONY: run
run:
	docker build -t golang-traning:week2 .
	docker run --rm --name golang-traning-week2 -p 8080:8080 golang-traning:week2

.PHONY: k6-test
k6-test:
	docker run -i loadimpact/k6 run - <week2/script.js

.PHONY: week2
week2:
	cd ./docker/k6-docker && \
	docker-compose up -d influxdb grafana goapp

.PHONY: test
test:
	cd ./docker/k6-docker && \
	docker-compose run -v `pwd`/../../week2:/scripts k6 run /scripts/script.js

# database (mysql + gorm + mytop)
mysql:
	docker-compose down && \
	docker-compose up -d db
app:
	docker-compose up -d goapp mytop --build
goapp:
	docker exec -ti golang-traning-week3 sh

.PHONY: mysql app goapp
