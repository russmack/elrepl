elRepl
======

An elasticsearch repl written in Golang.


This project is in the earliest stage of development.

```
Commands:
alias
close
count
dir
doc
env
exit
flush
get
help
host
index
load
log
mapping
open
optimize
port
post
put
recovery
refresh
reindex
run
segments
settings
stats
status
version
duplicatescount
```

```
Example usage:

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

> alias
Request: http://10.1.1.12:9200/_aliases?pretty=true
{
  "movies-2014-05-04-2252" : {
    "aliases" : {
      "movies" : { }
    }
  }
}

> alias move localhost 9200 fromIndex toIndex aliasName
...

> reindex localhost 9200 srcindex type localhost 9200 targetindex routing
...

> duplicatescount localhost 9200 index type field
...
```
