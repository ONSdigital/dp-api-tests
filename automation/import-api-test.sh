#!/bin/bash

docker-compose down

docker-compose run import_api_tests

docker-compose down