package commands

import (
	"github.com/russmack/elrepl/backend"
	"strconv"
)

func init() {
	h := backend.NewHandler()
	h.CommandName = "log"
	h.CommandPattern = "(log)( )(.*)"
	h.Usage = "log"
	h.HandlerFunc = func(cmd *backend.Command) (string, bool) {
		backend.LogLevel = 1
		return "Logging level set to: " + strconv.Itoa(backend.LogLevel), true
	}
	backend.HandlerRegistry[h.CommandName] = h
}
