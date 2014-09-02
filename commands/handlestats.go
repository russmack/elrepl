package commands

import (
	"fmt"
	"github.com/russmack/elrepl/backend"
	"net/url"
	"regexp"
	"strings"
)

type StatsCmd struct{}

func init() {
	h := backend.NewHandler()
	h.CommandName = "stats"
	h.CommandPattern = "(stats)( )(.*)"
	h.Usage = "stats indexName"
	h.CommandParser = func(cmd *backend.Command) (string, bool) {
		pattFn := map[*regexp.Regexp]func([]string) (string, bool){
			// Get stats
			regexp.MustCompile(`^stats ([a-zA-Z0-9\.\-]+)$`): func(s []string) (string, bool) {
				d := backend.Resource{
					Scheme:   "http",
					Host:     backend.Server.Host,
					Port:     backend.Server.Port,
					Index:    s[1],
					Endpoint: "_stats",
				}
				c := StatsCmd{}
				r, ok := c.Get(d)
				return r, ok
			},
			// Stats help
			regexp.MustCompile(`^stats /\?$`): func(s []string) (string, bool) {
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

func (c *StatsCmd) Get(d backend.Resource) (string, bool) {
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
