package main

import (
	"fmt"
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

		//fmt.Println("cmd:", cmd.Name)
		//arg := cmd.Args[0]
		//fmt.Println("arg:", arg)
		//fmt.Println("arg1:", cmd.Args[1])
		//fmt.Println("args:", cmd.Args)
		//bodyIdx := strings.Index(arg, "{")
		//fmt.Println("bodyidx:", bodyIdx)

		//queryArgs := arg[:bodyIdx]
		//queryArgs = strings.TrimPrefix(queryArgs, "/")
		//queryArgs = strings.TrimSpace(queryArgs)
		//body := arg[bodyIdx:]
		queryArgs := cmd.Args[0]
		body := cmd.Args[1]

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
