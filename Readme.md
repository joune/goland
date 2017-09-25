Adventures in "Goland" with Redis, ScyllaDB and Kafka
=====================================================

A service to categorize geolocated users. More here:
gist.github.com/daedric/db45c531a1bc5e58f0383f9c1bff4306

Modules
-------

* loader: loads users to redis and some scenarios to kafka
* consumer: listens to kafka 'sessions' topic to enrich and store them to scyllaDB
* server: a gRPC server that sorts a user's friends into categories
* client: a gRPC client that connects to the server and requests categories for its userId
* categories: functions implementing the friends categories for the server
* data: protobuf messages and services definitions
* geo: geoloc utilities
* util: sorting and grouping utilities
* db: connection, schema and queries to scyllaDB
* users: a repository of users stored in Redis

Prerequisites
-------------

* golang
* docker and docker-compose
* protobuf and protoc-gen-go

How to run
----------

* Run the stack

    > run.sh

* Try requests from client (x = 1 to 6)

    > ./client/client <x>

* When done, 

    > ./cleanup.sh

Expected results
----------------

* 1 and 2 are mutual love
* 3 and 4 are best friends
* 3 and 6 are crush
* 'empty' categories contain the client ID, don't be surprised!

Known issues
------------

* Kafka behaves strangely in docker (or maybe i'm doing something wrong); consumer needs to be restarted :(
