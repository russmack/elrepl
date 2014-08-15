package main

import (
	"fmt"
	"strings"
)

// curl -XPOST "http://localhost:9200/movies/_search?pretty" -d'{ ... body ... }''
// becomes
// post _search?pretty { "query": { "term": { "director": "scott" } } }
// Currently, must be on single line.
func init() {
	h := NewHandler()
	h.CommandName = "post"
	h.CommandPattern = "((?i)post(?-i))( )(.*)"
	h.HandlerFunc = func(cmd *Command) string {
		queryHost := server.host
		queryPort := server.port

		arg := cmd.Args[0]
		bodyIdx := strings.Index(arg, "{")
		queryArgs := arg[:bodyIdx]
		queryArgs = strings.TrimPrefix(queryArgs, "/")
		queryArgs = strings.TrimSpace(queryArgs)
		body := arg[bodyIdx:]

		url := fmt.Sprintf("http://%s:%s/%s", queryHost, queryPort, queryArgs)

		fmt.Println("Request:", url)
		res, err := postHttpResource(url, body)
		if err != nil {
			return err.Error()
		}
		return res
	}
	HandlerRegistry[h.CommandName] = h
}
