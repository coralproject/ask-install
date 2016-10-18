#!/bin/bash

docker-compose stop
docker-compose rm -f
docker volume rm $(docker volume ls -f dangling=true -q)
