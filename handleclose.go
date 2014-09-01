package main

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

type CloseCmd struct{}

func init() {
	h := NewHandler()
	h.CommandName = "close"
	h.CommandPattern = "(close)( )(.*)"
	h.Usage = "close indexName"
	h.CommandParser = func(cmd *Command) (string, bool) {
		pattFn := map[*regexp.Regexp]func([]string) (string, bool){
			// Do close
			regexp.MustCompile(`^close ([a-zA-Z0-9\.\-]+)$`): func(s []string) (string, bool) {
				d := Resource{
					Scheme:   "http",
					Host:     server.host,
					Port:     server.port,
					Index:    s[1],
					Endpoint: "_close",
				}
				c := CloseCmd{}
				r, ok := c.Do(d)
				return r, ok
			},
			// Close help
			regexp.MustCompile(`^close /\?$`): func(s []string) (string, bool) {
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

func (c *CloseCmd) Do(d Resource) (string, bool) {
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
