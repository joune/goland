Project
=======

gist.github.com/daedric/db45c531a1bc5e58f0383f9c1bff4306

Golang
======

* https://golang.org/doc/code.html
* https://tour.golang.org
* https://golang.org/pkg/
* https://github.com/fatih/vim-go

Env
===

Protobuf
--------

* install https://github.com/google/protobuf/releases/download/v3.4.0/protoc-3.4.0-linux-x86_64.zip
* go get -u github.com/golang/protobuf/protoc-gen-go
* protoc --go_out=. data/*.proto

Implem notes
============

Most seen - base query
---------

select user2, sum(duration) as t
  from sessions 
  where 
      user=client 
  and start_time >= (now - 7 days) 
  order by t desc limit 1

**InvalidRequest**: Error from server: code=2200 [Invalid query] message="the select clause must either contains only aggregates or none"

**InvalidRequest**: Error from server: code=2200 [Invalid query] message="Aliases are not allowed in order by clause ('t')"


Best friend
-----------

  'most seen' and location1=OTHER 

ServerError: Not implemented: INDEXES --> https://github.com/scylladb/scylla/issues/2025

When adding more PKs:

InvalidRequest: Error from server: code=2200 [Invalid query] message="Clustering column "is_night" cannot be restricted (preceding column "start_date" is restricted by a non-EQ relation)"


Crush
-----

 'most seen' and atNight=True and (location1=HOME or location2=HOME) && nbResults >= 3

Mutual Love
-----------

* use a counter to aggregate global time spent


Rich sessions
-------------

* location1, location2 = HOME | WORK | OTHER
* duration (start - end)
* atNight if starts before 8am 
          && ends after 10pm (j-1) 
          && min(end_time) - max(start_time) >= 6h

