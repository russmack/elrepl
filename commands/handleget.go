package commands

import (
	"fmt"
	"github.com/russmack/elrepl/backend"
	"net/url"
)

func init() {
	h := backend.NewHandler()
	h.CommandName = "get"
	h.CommandPattern = "(get)( )(.*)"
	h.Usage = "get url"
	h.HandlerFunc = func(cmd *backend.Command) (string, bool) {
		//queryHost := cmd.Tokens[1]
		//queryPort := cmd.Tokens[2]
		//queryArgs := cmd.Tokens[3]

		arg := cmd.Tokens[1]
		u := new(url.URL)
		newUrl, err := u.Parse(arg)
		if err != nil {
			return "Unable to parse url: " + err.Error(), false
		}
		fmt.Println("Request: get", newUrl)
		res, err := backend.GetHttpResource(newUrl.String())
		if err != nil {
			return err.Error(), false
		}
		return res, true
	}
	backend.HandlerRegistry[h.CommandName] = h
}
