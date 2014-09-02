package commands

import (
	"fmt"
	"github.com/russmack/elrepl/backend"
	"net/url"
	"regexp"
	"strings"
)

type SettingsCmd struct{}

func init() {
	h := backend.NewHandler()
	h.CommandName = "settings"
	h.CommandPattern = "(settings)( )(.*)"
	h.Usage = "settings indexName"
	h.CommandParser = func(cmd *backend.Command) (string, bool) {
		pattFn := map[*regexp.Regexp]func([]string) (string, bool){
			// Get settings
			regexp.MustCompile(`^settings ([a-zA-Z0-9\.\-]+)$`): func(s []string) (string, bool) {
				d := backend.Resource{
					Scheme:   "http",
					Host:     backend.Server.Host,
					Port:     backend.Server.Port,
					Index:    s[1],
					Endpoint: "_settings",
				}
				c := SettingsCmd{}
				r, ok := c.Get(d)
				return r, ok
			},
			// Settings help
			regexp.MustCompile(`^settings /\?$`): func(s []string) (string, bool) {
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

func (c *SettingsCmd) Get(d backend.Resource) (string, bool) {
	u := new(url.URL)
	u.Scheme = d.Scheme
	u.Host = d.Host + ":" + d.Port
	index := d.Index
	u.Path = index + "/" + d.Endpoint
	q := u.Query()
	q.Add("pretty", "true")
	u.RawQuery = q.Encode()
	fmt.Println("Request:", u)
	res, err := backend.GetHttpResource(u.String())
	if err != nil {
		return err.Error(), false
	}
	return res, true
}
