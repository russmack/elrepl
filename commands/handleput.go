package commands

import (
	"fmt"
	"github.com/russmack/elrepl/backend"
)

// curl -XPUT "http://localhost:9200/movies/movie/1" -d'{ ... body ... }''
// becomes
// put movie/1 { "title": "Alien", "director": "Ridley Scott", "year": 1979, "genres": ["Science fiction"] }
// TODO: Currently, must be on single line.
func init() {
	h := backend.NewHandler()
	h.CommandName = "put"
	h.CommandPattern = "(put)( )(.*)"
	h.Usage = `put movie/1 { "title": "Alien", "director": "Ridley Scott", "year": 1979, "genres": ["Science fiction"] }`
	h.HandlerFunc = func(cmd *backend.Command) (string, bool) {
		queryHost := cmd.Tokens[1]
		queryPort := cmd.Tokens[2]
		queryArgs := cmd.Tokens[3]
		url := fmt.Sprintf("http://%s:%s/%s", queryHost, queryPort, queryArgs)
		//url := fmt.Sprintf("http://%s:%s/%s/%s", queryHost, queryPort, backend.Server.Index, queryArgs)
		fmt.Println("Request: put", url)
		res, err := backend.PutHttpResource(url, cmd.Body)
		if err != nil {
			return err.Error(), false
		}
		return res, true
	}
	backend.HandlerRegistry[h.CommandName] = h
}
