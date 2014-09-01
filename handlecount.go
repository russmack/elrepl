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
	h.Usage = "count [host port] indexName [type]"
	h.CommandParser = func(cmd *Command) (string, bool) {
		pattFn := map[*regexp.Regexp]func([]string) (string, bool){
			// Get count: require: index; optional: type
			regexp.MustCompile(`^count ([a-zA-Z0-9\.\-]+)(( )([a-zA-Z0-9\.\-]+))?$`): func(s []string) (string, bool) {
				d := Resource{
					Scheme:   "http",
					Host:     server.host,
					Port:     server.port,
					Index:    s[1],
					Type:     s[4],
					Endpoint: "_count",
				}
				c := CountCmd{}
				r, ok := c.Get(d)
				return r, ok
			},
			// Get count: require: host, port, index; optional: type
			regexp.MustCompile(`^count ([a-zA-Z0-9\.\-]+) ([0-9]{1,5}) ([a-zA-Z0-9\.\-]+)(( )([a-zA-Z0-9\.\-]+))?$`): func(s []string) (string, bool) {
				d := Resource{
					Scheme:   "http",
					Host:     s[1],
					Port:     s[2],
					Index:    s[3],
					Type:     s[6],
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
	if d.Host == "" || d.Port == "" {
		return "Missing host or port.", false
	}
	u := new(url.URL)
	u.Scheme = d.Scheme
	u.Host = d.Host + ":" + d.Port
	index := d.Index + "/"
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
