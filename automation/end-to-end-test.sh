#!/bin/bash

docker-compose down

docker-compose up -d vault

docker-compose up -d kafka

sleep 5

export VAULT_TOKEN=$(docker-compose logs vault | awk '/Root Token:/ {print $NF;}' | tail -n 1 | sed 's/\x1b\[[0-9;]*m//g')

echo "acceptance test vault token: ${VAULT_TOKEN}"

docker-compose run end_to_end_tests

docker-compose down