package handlers

import (
	"fmt"
	"github.com/russmack/elrepl/types"
)

func init() {
	fmt.Println("REGISTERING HANDLER: VERSION")
	h := types.NewHandler()
	h.CommandName = "version"
	h.CommandPattern = "(version)(( )(.*))"
	h.HandlerFunc = func(cmd *types.Command) string {
		return `
			elRepl version 0.1
		`
	}
	HandlerRegistry[h.CommandName] = h
}
