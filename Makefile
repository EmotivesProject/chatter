SLEEP_TIME=0

.PHONY: lint test integration

lint:
	golangci-lint run

test:
	go test -v ./...

integration:
	make start_test_containers
	sleep $(SLEEP_TIME)
	go test -v -tags=integration ./...
	docker-compose -f docker/chatter/docker-compose.test.yml down

start_test_containers:
	make get_latest_containers
	docker-compose -f docker/chatter/docker-compose.test.yml down --remove-orphans
	docker-compose -f docker/chatter/docker-compose.test.yml up -d --remove-orphans

stop_test_containers:
	docker-compose -f docker/chatter/docker-compose.test.yml down --remove-orphans

integration_dry:
	go test -v -tags=integration ./...

get_latest_containers:
	docker pull ghcr.io/tombowyerresearchproject/postgres_db:latest
	docker pull ghcr.io/tombowyerresearchproject/uacl_api:latest

destory_test_containers:
	docker-compose -f docker/chatter/docker-compose.test.yml down --remove-orphans