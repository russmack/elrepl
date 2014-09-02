package commands

import (
	"fmt"
	"github.com/russmack/elrepl/backend"
	"net/url"
	"regexp"
	"strings"
)

type VersionCmd struct{}

func init() {
	h := backend.NewHandler()
	h.CommandName = "version"
	h.CommandPattern = "(version)(( )(.*))"
	h.Usage = "version host port"
	h.CommandParser = func(cmd *backend.Command) (string, bool) {
		pattFn := map[*regexp.Regexp]func([]string) (string, bool){
			// Get version
			regexp.MustCompile(`^version ([a-zA-Z0-9\.\-]+) ([0-9]{1,5})$`): func(s []string) (string, bool) {
				d := backend.Resource{
					Scheme: "http",
					Host:   s[1],
					Port:   s[2],
				}
				c := VersionCmd{}
				r, ok := c.Get(d)
				return r, ok
			},
			// Version help
			regexp.MustCompile(`^version /\?$`): func(s []string) (string, bool) {
				return "", false
			},
		}
		r, ok := h.Tokenize(strings.TrimSpace(cmd.Instruction), pattFn)
		return r, ok
	}
	h.HandlerFunc = func(cmd *backend.Command) (string, bool) {
		r, ok := h.CommandParser(cmd)
		if !ok {
			if r != "" {
				r += "\n\n"
			}
			return h.Usage, false
		}
		return r, true
	}
	backend.HandlerRegistry[h.CommandName] = h
}

func (c *VersionCmd) Get(d backend.Resource) (string, bool) {
	u := new(url.URL)
	u.Scheme = d.Scheme
	u.Host = d.Host + ":" + d.Port
	fmt.Println("Request:", u)
	res, err := backend.GetHttpResource(u.String())
	if err != nil {
		return err.Error(), false
	}
	return res, true
}
