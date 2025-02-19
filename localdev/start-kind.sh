#!/bin/bash

cd "$(dirname "${BASH_SOURCE[0]}")"

config=$1
KIND_EXPERIMENTAL_PROVIDER=podman kind create cluster --image kindest/node:v1.29.2 --config=$1
kubectl create configmap awscreds --from-env-file=../local-secret.env
./deployredis.sh
./deploypg.sh
./deploynginx.sh
