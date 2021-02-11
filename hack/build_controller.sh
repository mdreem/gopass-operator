#!/bin/bash

eval $(minikube docker-env)
cd ../controller
IMG="gopass-controller:latest" make docker-build
