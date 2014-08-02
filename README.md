elRepl
======

An elasticsearch repl written in Golang.


This project is in the earliest stage of development.

Commands:
version
help
exit
log
dir
load
run
host
port
get
put
post
reindex
duplicatescount

Example usage:

```
> host localhost
Set server host: localhost

> host
Server host: localhost

> port 9200
Set server port: 9200

> port
Server port: 9200

> index podcasts
Set index: podcasts

> get _aliases
Request: http://localhost:9200/movies/_aliases
{"movies-2014-05-04-2252":{"aliases":{"movies":{}}}}

> get _search?q=title:thx1138
Request: http://localhost:9200/movies/_search?q=title:thx1138
{"took":5,"timed_out":false,"_shards":{"total":5,"successful":5....

> reindex localhost:9200/srcindex/type localhost:9200/targetindex/routing
...

> duplicatescount localhost:9200/podcasts/channel/title
...
```
