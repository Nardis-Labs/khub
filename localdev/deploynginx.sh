#!/bin/bash

## validate correct kube-context
kube_context=$(kubectl config current-context)
if [ $kube_context = "kind-kind" ]; then
    kubectl apply --filename https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/static/provider/kind/deploy.yaml
    kubectl wait --namespace ingress-nginx \
    --for=condition=ready pod \
    --selector=app.kubernetes.io/component=controller \
    --timeout=120s
    exit 0
else
  echo "To avoid damage to real kubernetes clusters, this script has failed because your cluster context is not set to point at the local kind cluster."
  exit 1
fi