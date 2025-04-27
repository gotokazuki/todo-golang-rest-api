.PHONY: run, build, push, inspect, generate_oas, run_oas

default:
	cat Makefile

run: build
	docker run -p 8080:8080 todo-golang-rest-api

build:
	go mod tidy
	docker build -t todo-golang-rest-api .

push: build
	@if [ -z "$(VERSION)" ]; then \
		echo "Error: VERSION is required. Usage: make push VERSION=v1.0.0"; \
		exit 1; \
	fi
	@if [ -z "$(ECR_IMAGE_NAME)" ]; then \
		echo "Error: ECR_IMAGE_NAME is required. Usage: make push ECR_IMAGE_NAME=todo-golang-rest-api"; \
		exit 1; \
	fi
	aws ecr get-login-password --region ap-northeast-1 | docker login --username AWS --password-stdin $(shell aws sts get-caller-identity --query Account --output text).dkr.ecr.ap-northeast-1.amazonaws.com
	docker tag todo-golang-rest-api:latest $(shell aws sts get-caller-identity --query Account --output text).dkr.ecr.ap-northeast-1.amazonaws.com/$(ECR_IMAGE_NAME):$(VERSION)
	docker push $(shell aws sts get-caller-identity --query Account --output text).dkr.ecr.ap-northeast-1.amazonaws.com/$(ECR_IMAGE_NAME):$(VERSION)

inspect:
	docker inspect todo-golang-rest-api:latest

generate_oas:
	npx @redocly/cli build-docs docs/openapi.yaml -o docs/index.html

run_oas:
	docker run -p 8081:8080 -e SWAGGER_JSON=/todo/openapi.yaml -v ${PWD}/docs:/todo swaggerapi/swagger-ui
