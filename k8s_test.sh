#!/bin/bash

kubectl -n services port-forward svc/postgres-postgresql 5432:5432 & PF1=$!
kubectl -n services port-forward svc/broker 5672:5672 & PF2=$!

NS=villas-demo

POSTGRES_USER=$(kubectl -n ${NS} get secret postgres-credentials -o json | jq -r .data.username | base64 -d)
POSTGRES_PASS=$(kubectl -n ${NS} get secret postgres-credentials -o json | jq -r .data.password | base64 -d)
RABBITMQ_USER=$(kubectl -n ${NS} get secret postgres-credentials -o json | jq -r .data.username | base64 -d)
RABBITMQ_PASS=$(kubectl -n ${NS} get secret postgres-credentials -o json | jq -r .data.password | base64 -d)

go test ./routes/healthz -p 1 --args \
	-mode test \
	-db-host localhost -db-name villas \
	-db-user ${POSTGRES_USER} -db-pass ${POSTGRES_PASS} \
	-amqp-host localhost \
	-amqp-user ${RABBITMQ_USER} -amqp-pass ${RABBITMQ_PASS}

kill $PF1 $PF2
wait $PF1 $PF2
