package commands

import (
	"github.com/russmack/elrepl/backend"
	"sort"
)

func init() {
	h := backend.NewHandler()
	h.CommandName = "help"
	h.CommandPattern = "(help)( )(.*)"
	h.Usage = "help"
	h.HandlerFunc = func(cmd *backend.Command) (string, bool) {
		cmdList := ""
		keys := []string{}
		for k, _ := range backend.HandlerRegistry {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, s := range keys {
			//cmdList += s + ": " + backend.HandlerRegistry[s].Usage + "\n"
			cmdList += s + "\n"
		}
		return `
	Help
	----


Write commands require host and port be specified to avoid slip-ups.
Read commands can omit host and port, and session settings will be used.


Commands
--------
` + cmdList + `
`, true
	}
	backend.HandlerRegistry[h.CommandName] = h
}
