#!/bin/bash

docker-compose down

docker-compose up -d vault

export VAULT_TOKEN=$(docker-compose logs vault | awk '/Root Token:/ {print $NF;}' | tail -n 1 | sed 's/\x1b\[[0-9;]*m//g')

echo "acceptance test vault token: ${VAULT_TOKEN}"

docker-compose run download_service_web_tests

docker-compose down