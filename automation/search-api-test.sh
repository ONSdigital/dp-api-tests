#!/bin/bash

docker-compose down

docker-compose run search_api_tests

docker-compose down