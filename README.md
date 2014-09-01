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
index
load
log
mapping
open
optimize
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

> env host localhost
Host set: localhost

> env port 9200
Port set: 9200

> env index MyIndex
Index set: MyIndex

> env

        elRepl version 0.1

        Host: localhost
        Port: 9200
        Index: MyIndex

> alias
Request: http://localhost:9200/_aliases?pretty=true
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
