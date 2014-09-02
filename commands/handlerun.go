package commands

import (
	"github.com/russmack/elrepl/backend"
	"strings"
)

func init() {
	h := backend.NewHandler()
	h.CommandName = "run"
	h.CommandPattern = "(run)( )(.*)"
	h.Usage = "run"
	h.HandlerFunc = func(cmd *backend.Command) (string, bool) {
		verb := strings.ToLower(cmd.Tokens[0])
		cmdParser := backend.NewCommandParser()
		newCmd, err := cmdParser.Parse(backend.LoadedRequest.Request)
		if err != nil {
			return "Unable to parse loaded query for run command.", false
		}
		dispatcher := backend.NewDispatcher()
		switch verb {
		case "post":
			resp, ok := dispatcher.Dispatch(newCmd)
			return resp, ok
		case "put":
			resp, ok := dispatcher.Dispatch(newCmd)
			return resp, ok
		case "get":
			resp, ok := dispatcher.Dispatch(newCmd)
			return resp, ok
		}
		return "Unable to run loaded query.", false
	}
	backend.HandlerRegistry[h.CommandName] = h
}
