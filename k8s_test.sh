#!/bin/bash

kubectl -n services port-forward svc/postgres-postgresql 5432:5432 & PF=$!

POSTGRES_PASS=$(kubectl -n villas-demo get secret postgres-credentials -o json | jq -r .data.password | base64 -d)

go run start.go -dbhost localhost -dbname villas -dbuser villas -dbpass ${POSTGRES_PASS}

kill $PF
wait $PF
