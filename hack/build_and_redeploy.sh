#!/bin/bash

cd ..
cd gopass-server
make docker-login
make docker-build
make docker-push

cd ..
cd controller
make docker-login
make docker-build
make docker-push

kubectl rollout restart deployment operator-gopass-repository-deployment -n operator-system
kubectl rollout restart deployment operator-controller-manager -n operator-system
