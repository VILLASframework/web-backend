#!/usr/bin/env bash


cd ../../
go mod tidy
swag init -p pascalcase -g "start.go" -o "./doc/autoapi/"
cd -

redoc-cli bundle --cdn --title "VILLASweb Backend API" --output index.html swagger.yaml