package commands

import (
	"github.com/russmack/elrepl/backend"
	"io/ioutil"
)

func init() {
	h := backend.NewHandler()
	h.CommandName = "dir"
	h.CommandPattern = "(dir)( )(.*)"
	h.Usage = "dir"
	h.HandlerFunc = func(cmd *backend.Command) (string, bool) {
		dir := ""
		if len(cmd.Tokens) == 1 {
			dir = "."
		} else {
			dir = cmd.Tokens[1]
		}
		dirFiles, err := ioutil.ReadDir(dir)
		if err != nil {
			return err.Error(), false
		}
		files := ""
		for _, j := range dirFiles {
			files += j.Name() + "\n"
		}
		return files, true
	}
	backend.HandlerRegistry[h.CommandName] = h
}
