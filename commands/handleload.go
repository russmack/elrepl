package commands

import (
	"github.com/russmack/elrepl/backend"
	"io/ioutil"
)

func init() {
	h := backend.NewHandler()
	h.CommandName = "load"
	h.CommandPattern = "(load)( )(.*)"
	h.Usage = "load filename"
	h.HandlerFunc = func(cmd *backend.Command) (string, bool) {
		if len(cmd.Tokens) != 2 {
			return h.Usage, false
		}
		filename := cmd.Tokens[1]
		file, err := ioutil.ReadFile(filename)
		if err != nil {
			return err.Error(), false
		}
		fileText := string(file)
		backend.LoadedRequest.Request = fileText
		return fileText, true
	}
	backend.HandlerRegistry[h.CommandName] = h
}
