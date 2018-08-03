#!/bin/bash

docker-compose down

docker-compose run search_api_web_tests

docker-compose down
