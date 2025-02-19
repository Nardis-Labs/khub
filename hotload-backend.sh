#!/bin/bash
pgrep kubectl port-forward | xargs kill

GOOS=linux go build -o ./tmp/khub .

DOCKER_BUILDKIT=1 podman build -f Dockerfile.localdev -t khublocal:0.0.1 .
podman save --format docker-archive khublocal:0.0.1 -o khublocal.tar
KIND_EXPERIMENTAL_PROVIDER=podman kind load image-archive ./khublocal.tar
kubectl delete -f ./localdev/khub.yaml
kubectl apply -f ./localdev/khub.yaml
sleep 2
kubectl port-forward svc/khub-app 8080:8080 &