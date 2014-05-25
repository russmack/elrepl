elrepl
======

An elasticsearch repl written in Golang.


Example usage:


> host localhost
Set server host: localhost

> port 9200
Set server port: 9200

> index podcasts
Set index: podcasts

> get _aliases
Request: http://localhost:9200/movies/_aliases
{"movies-2014-05-04-2252":{"aliases":{"movies":{}}}}

> get _search?q=title:thx1138
Request: http://localhost:9200/movies/_search?q=title:thx1138
{"took":5,"timed_out":false,"_shards":{"total":5,"successful":5....

