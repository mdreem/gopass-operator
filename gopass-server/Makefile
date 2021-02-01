GOPATH:=$(shell go env GOPATH)
GOMOD =$(shell go env GOMOD)
PROTOC_GEN_GO := $(GOPATH)/bin/protoc-gen-go
SHELL := /bin/bash
PATH := $(GOPATH)/bin/:$(PATH)

$(PROTOC_GEN_GO):
	go get -u github.com/golang/protobuf/protoc-gen-go

%.pb.go: %.proto | $(PROTOC_GEN_GO)
	protoc --go_out=plugins=grpc:. $<

build_protobuf: gopass-server/gopass_repository/repository.pb.go

build: build_protobuf
	go build -o bin/gopass_server cmd/main.go

build_docker: build
	docker build . -f gopass-server/docker/Dockerfile -t gopass-server
