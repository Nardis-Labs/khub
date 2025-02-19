#!/bin/bash

helm install postgres --set auth.postgresPassword=postgres1011 oci://registry-1.docker.io/bitnamicharts/postgresql -f ./pg_values.yaml