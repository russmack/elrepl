package handlers

import (
	"fmt"
	"github.com/russmack/elrepl/types"
	"os"
)

func init() {
	fmt.Println("REGISTERING HANDLER: VERSION")
	h := types.NewHandler()
	h.CommandName = "exit"
	h.CommandPattern = "(exit)( )(.*)"
	h.HandlerFunc = func(cmd *types.Command) string {
		fmt.Println("Bye.")
		os.Exit(0)
		return ""
	}
	HandlerRegistry[h.CommandName] = h
}
