#!/bin/bash

docker-compose down

docker-compose run filter_api_web_tests

docker-compose down