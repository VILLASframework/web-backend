#!/bin/bash

set -e

NS=${NS:-villas}
MK="minikube"
KCTL="kubectl -n${NS}"
HELM="helm -n${NS}"
# CHART="villas"
CHART="~/workspace/rwth/acs/public/catalogue/charts/villas"

CONFIG=$(mktemp)
NAME="villas"

function cleanup() {
	kill $PF1 $PF2 $PF3
	wait $PF1 $PF2 $PF3

	rm ${CONFIG}

	echo "Goodbye"
}
trap cleanup EXIT

if [ -n "${USE_MINIKUBE}" ]; then
	MK_START_OPTS="--addons=ingress"
	if [ $(uname -s) == Darwin ]; then
		# https://github.com/kubernetes/minikube/issues/7332
		MK_START_OPTS+="--vm=true"
	fi

	${MK} start ${MK_START_OPTS}

	kubectl -n kube-system expose deployment ingress-nginx-controller --type=LoadBalancer || true

	IP=$(minikube ip)
	PORT=$(kubectl -n kube-system get service ingress-nginx-controller --output='jsonpath={.spec.ports[0].nodePort}')
	
	# Add pseudo hostname to /etc/hosts
	echo "Please provide your root password for modifiying /etc/hosts"
	sudo sed -in "/^${IP}.*/d" /etc/hosts
	echo "${IP} minikube" | sudo tee -a /etc/hosts

	HELM_OPTS="--set web.backend.enabled=false \
			   --set web.backend.external.ip=${IP} \
			   --set web.backend.external.port=4000 \
			   --set web.admin.password=admin \
			   --set ingress.host=minikube"

	# Check if chart has already been deployed before
	if helm get values villas > /dev/null; then
		RABBITMQ_ERLANG_COOKIE=$(kubectl get secret --namespace default villas-broker -o jsonpath="{.data.rabbitmq-erlang-cookie}" | base64 --decode)
		RABBITMQ_PASSWORD=$(kubectl get secret --namespace default villas-broker -o jsonpath="{.data.rabbitmq-password}" | base64 --decode)
		
		HELM_OPTS+=" --set broker.auth.password=${RABBITMQ_PASSWORD} --set broker.auth.erlangCookie=${RABBITMQ_ERLANG_COOKIE}"
	fi

	${HELM} upgrade \
		${HELM_OPTS} \
		--install \
		--create-namespace \
		--wait \
		--repo https://packages.fein-aachen.org/helm/charts \
		${NAME} ${CHART}
fi

# Get backend config from cluster
${KCTL} get cm ${NAME}-web -o 'jsonpath={.data.config\.yaml}' > ${CONFIG}

# Enable access of backend to db, broker and s3
export DB_HOST=localhost
export DB_PASS=$(${KCTL} get secret ${NAME}-database -o 'jsonpath={.data.postgresql-password}' | base64 -d)

export AMQP_HOST=localhost
export AMQP_PASS=$(${KCTL} get secret ${NAME}-broker -o 'jsonpath={.data.rabbitmq-password}' | base64 -d)

export AWS_ACCESS_KEY_ID=$(${KCTL} get secret ${NAME}-s3 -o 'jsonpath={.data.accesskey}' | base64 -d)
export AWS_SECRET_ACCESS_KEY=$(${KCTL} get secret ${NAME}-s3 -o 'jsonpath={.data.secretkey}' | base64 -d)

# Setup port forwards for backend
if [ -n "${USE_EXTERNAL_POSTGRESQL}" ]; then
	kubectl -n postgresql port-forward svc/postgresql 5432:5432 & PF1=$!
else
	${KCTL} port-forward svc/${NAME}-database 5432:5432 & PF1=$!
fi
${KCTL} port-forward svc/${NAME}-broker 5672:5672 & PF2=$!

sleep 2

if [ -n "${USE_MINIKUBE}" ]; then
	# python -mwebbrowser http://minikube:${PORT}
	echo
	echo "==========================================="
	echo "Access VILLASweb at http://minikube:${PORT}"
	echo "==========================================="
	echo
fi

go run start.go -config ${CONFIG}
