.PHONY: migrate test

include .env

DATABASE_TEST_HOST=boilerplate-database-test
DATABASE_TEST_PORT=5436

DATABASE_URL=postgres://${DATABASE_USER}:${DATABASE_PASSWORD}@localhost:${DATABASE_PORT}/${DATABASE_NAME}?sslmode=${DATABASE_SSL}

IMAGE=diegodesousas/boilerplate
SERVICE_NAME=boilerplate
BASEDIR=/application/boilerplate
DATABASE_VOLUME=boilerplate-db
DATABASE_TEST_VOLUME=boilerplate-db-test

NETWORK=boilerplate-net

UID=$(shell id -u ${USER})
GID=$(shell id -g ${USER})

build-development-image:
	@docker rmi -f ${IMAGE}-development
	@echo "Building image ${IMAGE}-development"
	@docker \
		build \
		--target development \
		-t ${IMAGE}-development \
		.

build-production-image:
	@docker rmi -f ${IMAGE}-production
	@echo "Building image ${IMAGE}-production"
	@docker \
		build \
		--target production \
		-t ${IMAGE}-production \
		.

mod:
	@docker run \
		-it \
		-v ${PWD}:${BASEDIR} \
		-p ${HTTP_PORT}:${HTTP_PORT} \
		--name ${SERVICE_NAME} \
		--rm  \
		-w ${BASEDIR} \
		--network ${NETWORK} \
		${IMAGE}-development \
		go mod vendor

dev-server: database-up
	@docker run \
		-it \
		-v ${PWD}:${BASEDIR} \
		-p ${HTTP_PORT}:${HTTP_PORT} \
		--name ${SERVICE_NAME} \
		--rm  \
		-w ${BASEDIR} \
		--network ${NETWORK} \
		${IMAGE}-development \
		go run entrypoints/server/main.go

test:
	@docker run \
		-it \
		-v ${PWD}:${BASEDIR} \
		-p ${HTTP_PORT}:${HTTP_PORT} \
		--name ${SERVICE_NAME} \
		--rm  \
		-w ${BASEDIR} \
		--network ${NETWORK} \
		-e BASE_PROJECT_PATH=${BASEDIR} \
		${IMAGE}-development \
		/bin/bash -c "go test ./..."

build-network:
ifeq ($(shell docker network list --filter name=${NETWORK} | wc -l), 1)
	@echo "Creating network ${NETWORK}"
	@docker network create --driver bridge ${NETWORK}
else
	@echo "Network ${NETWORK} created"
endif

build-volume:
ifeq ($(shell docker volume ls -f name=${DATABASE_VOLUME} | wc -l), 1)
	@echo "Creating db volume ${DATABASE_VOLUME}"
	@docker volume create ${DATABASE_VOLUME}
else
	@echo "Volume ${DATABASE_VOLUME} created"
endif

database-up:
ifeq ($(shell docker ps -f name=^/${DATABASE_HOST}$$ | wc -l), 1)
	@echo "Starting ${DATABASE_HOST} database service"
	@docker run \
		--rm \
		-v ${DATABASE_VOLUME}:/var/lib/postgresql/data \
		-p ${DATABASE_PORT_LOCAL}:${DATABASE_PORT} \
		--name ${DATABASE_HOST} \
		--hostname ${DATABASE_HOST} \
		-e POSTGRES_USER=${DATABASE_USER} \
		-e POSTGRES_PASSWORD=${DATABASE_PASSWORD} \
		-e POSTGRES_DB=${DATABASE_NAME} \
		--network ${NETWORK} \
		-d \
		postgres:12
else
	 @echo "Database ${DATABASE_HOST} up and running"
endif

database-down:
	@echo "stoping $(DATABASE_HOST)"
	@docker stop $(DATABASE_HOST)
	@echo "postgres $(DATABASE_HOST) stoped."

migrate-create:
	@migrate create -ext sql -dir infra/database/migrations create_$(name)
	@echo "after the file be edited run: 'make migrate ACTION=up'"

migrate-update:
	@migrate create -ext sql -dir infra/database/migrations update_$(name)
	@echo "after the file be edited run: 'make migrate ACTION=up'"

# make migrate ACTION=up or make migrate ACTION=down
migrate:
	./migrate -path infra/database/migrations -database "$(DATABASE_URL)" $(ACTION)

dev-init: build-development-image build-network build-volume
