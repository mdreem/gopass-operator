SHELL := /bin/bash

build-server:
	./build_server.sh

build-controller:
	./build_controller.sh

all: build-server build-controller

restart-repository:
	kubectl rollout restart deployment operator-gopass-repository-deployment -n operator-system

restart-controller:
	kubectl rollout restart deployment operator-controller-manager -n operator-system

restart-all: restart-controller

build-restart: all restart-all

logs-server:
	kubectl logs "$$(kubectl get pods -n operator-system | grep "gopassrepository-[^ ]*" -o)" -n operator-system -f

logs-controller:
	kubectl logs "$$(kubectl get pods -n operator-system | grep "operator-controller-manager[^ ]*" -o)" manager -n operator-system -f

show-pods:
	kubectl get pods -n operator-system

show-deployments:
	kubectl get deployments -n operator-system

show-gopass-repositories:
	kubectl get gopassrepository
