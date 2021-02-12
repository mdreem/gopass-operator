#!/bin/bash

eval $(minikube docker-env)
cd ../gopass-server
IMG="gopass-server:latest" make docker-build
