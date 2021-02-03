GOPATH:=$(shell go env GOPATH)
GOMOD =$(shell go env GOMOD)
PROTOC_GEN_GO := $(GOPATH)/bin/protoc-gen-go
SHELL := /bin/bash
PATH := $(GOPATH)/bin/:$(PATH)

VERSION ?= 0.0.2
IMG ?= public.ecr.aws/d7n7y2i7/gopass-operator:$(VERSION)

$(PROTOC_GEN_GO):
	go get -u github.com/golang/protobuf/protoc-gen-go

%.pb.go: %.proto | $(PROTOC_GEN_GO)
	protoc --go_out=plugins=grpc:. $<

build_protobuf: gopass_repository/repository.pb.go

build: build_protobuf
	go build -o bin/gopass_server cmd/main.go

docker-build: build
	docker build . -f docker/Dockerfile -t $(IMG)

docker-push:
	docker push ${IMG}

docker-login:
	aws ecr-public get-login-password --region us-east-1 | docker login --username AWS --password-stdin public.ecr.aws
