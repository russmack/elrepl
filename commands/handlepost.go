package commands

import (
	"fmt"
	"github.com/russmack/elrepl/backend"
)

// curl -XPOST "http://localhost:9200/movies/_search?pretty" -d'{ ... body ... }''
// becomes
// post _search?pretty { "query": { "term": { "director": "scott" } } }
// TODO: Currently, must be on single line.
func init() {
	h := backend.NewHandler()
	h.CommandName = "post"
	h.CommandPattern = "((?i)post(?-i))( )(.*)"
	h.HandlerFunc = func(cmd *backend.Command) (string, bool) {
		fmt.Println("...>", cmd.Tokens)

		queryHost := cmd.Tokens[1]
		queryPort := cmd.Tokens[2]
		queryArgs := cmd.Tokens[3]
		url := fmt.Sprintf("http://%s:%s/%s", queryHost, queryPort, queryArgs)
		fmt.Println("Request: post", url)
		fmt.Println("Request: post body", cmd.Body)
		res, err := backend.PostHttpResource(url, cmd.Body)
		if err != nil {
			return err.Error(), false
		}
		return res, true
	}
	backend.HandlerRegistry[h.CommandName] = h
}
