package main

import (
	"fmt"
)

// curl -XPUT "http://localhost:9200/movies/movie/1" -d'{ ... body ... }''
// becomes
// put movie/1 { "title": "Alien", "director": "Ridley Scott", "year": 1979, "genres": ["Science fiction"] }
// TODO: Currently, must be on single line.
func init() {
	h := NewHandler()
	h.CommandName = "put"
	h.CommandPattern = "(put)( )(.*)"
	h.Usage = `put movie/1 { "title": "Alien", "director": "Ridley Scott", "year": 1979, "genres": ["Science fiction"] }`
	h.HandlerFunc = func(cmd *Command) string {
		queryArgs := cmd.Tokens[2]
		url := fmt.Sprintf("http://%s:%s/%s/%s", server.host, server.port, server.index, queryArgs)
		fmt.Println("Request:", url)
		res, err := putHttpResource(url, cmd.Body)
		if err != nil {
			return err.Error()
		}
		return res
	}
	HandlerRegistry[h.CommandName] = h
}
