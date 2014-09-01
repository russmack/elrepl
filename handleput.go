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
		queryHost := cmd.Tokens[1]
		queryPort := cmd.Tokens[2]
		queryArgs := cmd.Tokens[3]
		url := fmt.Sprintf("http://%s:%s/%s", queryHost, queryPort, queryArgs)
		//url := fmt.Sprintf("http://%s:%s/%s/%s", queryHost, queryPort, server.index, queryArgs)
		fmt.Println("Request: put", url)
		res, err := putHttpResource(url, cmd.Body)
		if err != nil {
			return err.Error()
		}
		return res
	}
	HandlerRegistry[h.CommandName] = h
}
