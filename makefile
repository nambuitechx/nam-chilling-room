.PHONY: run down

POSTGRES_IMAGE=nam-chilling-room/postgres
SERVER_IMAGE=nam-chilling-room/server
TAG=1.0.0

build-server:
	docker build -t $(SERVER_IMAGE):$(TAG) ./server

build: build-server

run:
	docker compose up -d

down:
	docker compose down