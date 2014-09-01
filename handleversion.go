package main

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

type VersionCmd struct{}

func init() {
	h := NewHandler()
	h.CommandName = "version"
	h.CommandPattern = "(version)(( )(.*))"
	h.Usage = "version host port"
	h.CommandParser = func(cmd *Command) (string, bool) {
		pattFn := map[*regexp.Regexp]func([]string) (string, bool){
			// Get version
			regexp.MustCompile(`^version ([a-zA-Z0-9\.\-]+) ([0-9]{1,5})$`): func(s []string) (string, bool) {
				d := Resource{
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

func (c *VersionCmd) Get(d Resource) (string, bool) {
	u := new(url.URL)
	u.Scheme = d.Scheme
	u.Host = d.Host + ":" + d.Port
	fmt.Println("Request:", u)
	res, err := getHttpResource(u.String())
	if err != nil {
		return err.Error(), false
	}
	return res, true
}
