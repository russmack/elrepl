package commands

import (
	"fmt"
	"github.com/russmack/elrepl/backend"
	"net/url"
	"strings"
)

func init() {
	h := backend.NewHandler()
	h.CommandName = "delete"
	h.CommandPattern = "(delete)( )(.*)"
	// TODO: usage
	h.HandlerFunc = func(cmd *backend.Command) (string, bool) {
		//arg := cmd.Args[0]
		arg := strings.SplitN(cmd.Instruction, " ", 2)[1]
		u := new(url.URL)
		newUrl, err := u.Parse(arg)
		if err != nil {
			return "Unable to parse url: " + err.Error(), false
		}
		fmt.Println("Request:", newUrl)
		res, err := backend.DeleteHttpResource(newUrl.String())
		if err != nil {
			return err.Error(), false
		}
		return res, true
	}
	backend.HandlerRegistry[h.CommandName] = h
}
