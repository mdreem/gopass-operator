GOPATH:=$(shell go env GOPATH)
GOMOD =$(shell go env GOMOD)
PROTOC_GEN_GO := $(GOPATH)/bin/protoc-gen-go
SHELL := /bin/bash
PATH := $(GOPATH)/bin/:$(PATH)

$(PROTOC_GEN_GO):
	go get -u github.com/golang/protobuf/protoc-gen-go

%.pb.go: %.proto | $(PROTOC_GEN_GO)
	protoc --go_out=plugins=grpc:. $<

compile: gopass-server/gopass_repository/repository.pb.go
