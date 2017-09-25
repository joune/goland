#!/bin/bash

# compile protobufs if needed
protoc  --go_out=plugins=grpc:. data/*.proto  

# start scyllaDB (create keyspace)
./scylla.sh

# start the rest
docker-compose up --build -d

# sorry: horrible workaround to 'hide' the kafka issue :/
echo "Bear with me for a few seconds... "
sleep 10 
docker-compose restart consumer
docker logs zenly_consumer_1


echo "... We should be ready. Let's give it a try."
go get golang.org/x/net/context
go get google.golang.org/grpc
cd client
go build
./client 3
cd ..

echo "--------"
echo "Try from other clients using"
echo "./client/client <x> # x from 1 to 6"
