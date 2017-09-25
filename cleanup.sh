#!/bin/bash

docker-compose down
docker rm -vf scylladb
docker system prune
