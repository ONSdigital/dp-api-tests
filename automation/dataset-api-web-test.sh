#!/bin/bash

docker-compose down

docker-compose run dataset_api_web_tests

docker-compose down