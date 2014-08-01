package handlers

import (
	"fmt"
	"github.com/russmack/elrepl/types"
)

func init() {
	fmt.Println("REGISTERING HANDLER: VERSION")
	h := types.NewHandler()
	h.CommandName = "help"
	h.CommandPattern = "(help)( )(.*)"
	h.HandlerFunc = func(cmd *types.Command) string {
		return `
			Help
			----
			Commands:
			  eg:
			  host localhost
			  port 9200
			  index movies
			  get _search?q=title:thx1138
			`
	}
	HandlerRegistry[h.CommandName] = h
}
