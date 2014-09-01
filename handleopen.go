package main

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

type OpenCmd struct{}

func init() {
	h := NewHandler()
	h.CommandName = "open"
	h.CommandPattern = "(open)( )(.*)"
	h.Usage = "open indexName"
	h.CommandParser = func(cmd *Command) (string, bool) {
		pattFn := map[*regexp.Regexp]func([]string) (string, bool){
			// Do open
			regexp.MustCompile(`^open ([a-zA-Z0-9\.\-]+)$`): func(s []string) (string, bool) {
				d := Resource{
					Scheme:   "http",
					Host:     server.host,
					Port:     server.port,
					Index:    s[1],
					Endpoint: "_open",
				}
				c := OpenCmd{}
				r, ok := c.Do(d)
				return r, ok
			},
			// Open help
			regexp.MustCompile(`^open /\?$`): func(s []string) (string, bool) {
				return "", false
			},
		}
		r, ok := h.Tokenize(strings.TrimSpace(cmd.Instruction), pattFn)
		return r, ok
	}
	h.HandlerFunc = func(cmd *Command) string {
		r, ok := h.CommandParser(cmd)
		if !ok {
			if r != "" {
				r += "\n\n"
			}
			return r + usageMessage(h.Usage)
		}
		return r
	}
	HandlerRegistry[h.CommandName] = h
}

func (c *OpenCmd) Do(d Resource) (string, bool) {
	u := new(url.URL)
	u.Scheme = d.Scheme
	u.Host = d.Host + ":" + d.Port
	index := d.Index
	u.Path = index + "/" + d.Endpoint
	q := u.Query()
	q.Add("pretty", "true")
	u.RawQuery = q.Encode()
	fmt.Println("Request:", u)
	res, err := postHttpResource(u.String(), "")
	if err != nil {
		return err.Error(), false
	}
	return res, true
}
