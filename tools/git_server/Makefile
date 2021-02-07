SHELL := /bin/bash

VERSION ?= 0.0.8
IMG ?= public.ecr.aws/d7n7y2i7/git-test-server:$(VERSION)

docker-build:
	docker build . -f ./Dockerfile -t $(IMG)

docker-push:
	docker push ${IMG}

docker-login:
	aws ecr-public get-login-password --region us-east-1 | docker login --username AWS --password-stdin public.ecr.aws

deploy:
	kubectl apply -f git-server-deployment.yaml
