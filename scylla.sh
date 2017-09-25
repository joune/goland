#!/bin/bash

# all these sh files and custom networks are here only to workaround
# the fact that I can't initialize the keyspace automatically in scylladb docker :(

docker network create --attachable zenly_db

# customize scylla to support counters
docker build -f Dockerfile.scylladb -t zenly/scylla .
docker run --name scylladb --network="zenly_db" -d zenly/scylla

CQL="CREATE KEYSPACE zenly WITH REPLICATION = {'class': 'SimpleStrategy', 'replication_factor': 1};"
until docker exec scylladb /bin/sh -c "echo \"$CQL\" | cqlsh"; do
  echo "cqlsh: ScyllaDB is unavailable - retry later"
  sleep 10 #takes almost a minute
done 
echo "KEYSPACE zenly CREATED"
