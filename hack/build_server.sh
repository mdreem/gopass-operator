#!/bin/bash

eval $(minikube docker-env)
cd ../gopass-server
IMG="gopass-operator:latest" make docker-build
