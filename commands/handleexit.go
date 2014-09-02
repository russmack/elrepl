package commands

import (
	"fmt"
	"github.com/russmack/elrepl/backend"
	"os"
)

func init() {
	h := backend.NewHandler()
	h.CommandName = "exit"
	h.CommandPattern = "(exit)( )(.*)"
	h.Usage = "exit"
	h.HandlerFunc = func(cmd *backend.Command) (string, bool) {
		fmt.Println("Bye.")
		os.Exit(0)
		return "", true
	}
	backend.HandlerRegistry[h.CommandName] = h
}
