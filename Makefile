.PHONY: gen-src unit-test lint swagger docker-up docker-down test-docker-up test-docker-down docker-clear

gen-src:
	@go generate ./src/...

unit-test:
	@go test -v -short ./src/...

lint:
	@cd src && golangci-lint run --config ../.github/.golangci.yml && cd ..

swagger:
	@swag init -d src

docker-up:
	@docker-compose -f docker/development/docker-compose.yml up -d
	sleep 10
	# @chmod -R 777 docker/development

docker-down:
	# @chmod -R 777 docker/development
	@docker-compose -f docker/development/docker-compose.yml down
	# @chmod -R 777 docker/development

test-docker-up:
	@docker-compose -f docker/test/docker-compose.yml up -d
	sleep 10
	# @chmod -R 777 docker/test

test-docker-down:
	@docker-compose -f docker/test/docker-compose.yml down

docker-clear:
	@rm -rf docker/development/mysql-data
	@rm -rf docker/test/mysql-data
