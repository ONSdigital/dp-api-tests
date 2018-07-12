#!/bin/bash

docker-compose down

docker-compose run hierarchy_api_tests

docker-compose down