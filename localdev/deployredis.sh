#!/bin/bash

namespace=redis
kube_context=$(kubectl config current-context)

if [ $kube_context = "kind-kind" ]; then
    kubectl create namespace $namespace
    helm install redis -n redis --set auth.enabled=false --set replica.replicaCount=2 oci://registry-1.docker.io/bitnamicharts/redis 
    exit 0
else
  echo "To avoid damage to real kubernetes clusters, this script has failed because your cluster context is not set to point at the local kind cluster."
  exit 1
fi



