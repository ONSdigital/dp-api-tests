#!/bin/bash

docker-compose down

docker-compose run dataset_api_publishing_tests

docker-compose down