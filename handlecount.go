package main

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

type CountCmd struct{}

func init() {
	h := NewHandler()
	h.CommandName = "count"
	h.CommandPattern = "(count)( )(.*)"
	h.Usage = "count [/|indexName] [type]"
	h.CommandParser = func(cmd *Command) (string, bool) {
		pattFn := map[*regexp.Regexp]func([]string) (string, bool){
			// Get count
			regexp.MustCompile(`^count ([a-zA-Z0-9\.\-]+) ([a-zA-Z0-9\.\-]+)$`): func(s []string) (string, bool) {
				d := Resource{
					Scheme:   "http",
					Host:     server.host,
					Port:     server.port,
					Index:    s[1],
					Type:     s[2],
					Endpoint: "_count",
				}
				c := CountCmd{}
				r, ok := c.Get(d)
				return r, ok
			},
			// Count help
			regexp.MustCompile(`^count /\?$`): func(s []string) (string, bool) {
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

func (c *CountCmd) Get(d Resource) (string, bool) {
	u := new(url.URL)
	u.Scheme = d.Scheme
	u.Host = d.Host + ":" + d.Port
	index := d.Index
	docType := d.Type
	if docType != "" {
		index += docType + "/"
	}
	u.Path = index + d.Endpoint
	q := u.Query()
	q.Add("pretty", "true")
	u.RawQuery = q.Encode()
	fmt.Println("Request:", u)
	res, err := getHttpResource(u.String())
	if err != nil {
		return err.Error(), false
	}
	return res, true
}
