package main

import (
	"fmt"
)

// curl -XPOST "http://localhost:9200/movies/_search?pretty" -d'{ ... body ... }''
// becomes
// post _search?pretty { "query": { "term": { "director": "scott" } } }
// TODO: Currently, must be on single line.
func init() {
	h := NewHandler()
	h.CommandName = "post"
	h.CommandPattern = "((?i)post(?-i))( )(.*)"
	h.HandlerFunc = func(cmd *Command) string {
		queryHost := cmd.Tokens[1]
		queryPort := cmd.Tokens[2]
		queryArgs := cmd.Tokens[3]
		url := fmt.Sprintf("http://%s:%s/%s", queryHost, queryPort, queryArgs)
		fmt.Println("Request: post", url)
		res, err := postHttpResource(url, cmd.Body)
		if err != nil {
			return err.Error()
		}
		return res
	}
	HandlerRegistry[h.CommandName] = h
}
