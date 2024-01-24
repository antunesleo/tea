build:
	docker build . --tag antunesleo/tea:latest

push:
	docker push antunesleo/tea:latest

run:
	docker-compose up