package main

import (
	"fmt"
	"strings"
)

// curl -XPUT "http://localhost:9200/movies/movie/1" -d'{ ... body ... }''
// becomes
// put movie/1 { "title": "Alien", "director": "Ridley Scott", "year": 1979, "genres": ["Science fiction"] }
// Currently, must be on single line.
func init() {
	h := NewHandler()
	h.CommandName = "put"
	h.CommandPattern = "(put)( )(.*)"
	h.HandlerFunc = func(cmd *Command) string {
		arg := cmd.Args[0]

		bodyIdx := strings.Index(arg, " ")
		queryArgs := ""
		if bodyIdx > -1 {
			queryArgs = arg[:bodyIdx]
		}
		body := ""
		if bodyIdx > -1 {
			body = arg[bodyIdx:]
		}

		url := fmt.Sprintf("http://%s:%s/%s/%s", server.host, server.port, server.index, queryArgs)
		fmt.Println("Request:", url)
		res, err := putHttpResource(url, body)
		if err != nil {
			return err.Error()
		}
		return res
	}
	HandlerRegistry[h.CommandName] = h
}
