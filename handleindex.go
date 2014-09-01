package main

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

type IndexCmd struct{}

func init() {
	h := NewHandler()
	h.CommandName = "index"
	h.CommandPattern = "(index)(( )(.*))"
	h.Usage = "index ( (create host port indexName) | (delete host port indexName) )"
	h.CommandParser = func(cmd *Command) (string, bool) {
		pattFn := map[*regexp.Regexp]func([]string) (string, bool){
			// Index help
			regexp.MustCompile(`^index /\?$`): func(s []string) (string, bool) {
				return "", false
			},
			// Create index
			//curl -XPUT "http://localhost:9200/indexname/"
			regexp.MustCompile(`^index create ([a-zA-Z0-9\.\-]+) ([0-9]{1,5}) ([a-zA-Z0-9\.\-]+)$`): func(s []string) (string, bool) {
				d := Resource{
					Scheme: "http",
					Host:   s[1],
					Port:   s[2],
					Index:  s[3],
				}
				c := IndexCmd{}
				r, ok := c.Create(d)
				return r, ok
			},
			// Delete index
			regexp.MustCompile(`^index delete ([a-zA-Z0-9\.\-]+) ([0-9]{1,5}) ([a-zA-Z0-9\.\-]+)$`): func(s []string) (string, bool) {
				d := Resource{
					Scheme: "http",
					Host:   s[1],
					Port:   s[2],
					Index:  s[3],
				}
				c := IndexCmd{}
				r, ok := c.Delete(d)
				return r, ok
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

func (c *IndexCmd) Create(d Resource) (string, bool) {
	//curl -XPUT "http://localhost:9200/indexname/"
	u := new(url.URL)
	u.Scheme = d.Scheme
	u.Host = d.Host + ":" + d.Port
	u.Path = d.Index
	q := u.Query()
	q.Add("pretty", "true")
	u.RawQuery = q.Encode()
	fmt.Println("Request:", u)
	err := errors.New("")
	res, err := putHttpResource(u.String(), "")
	if err != nil {
		return err.Error(), false
	}
	return res, true
}

func (c *IndexCmd) Delete(d Resource) (string, bool) {
	//curl -XDELETE "http://localhost:9200/indexname/"
	u := new(url.URL)
	u.Scheme = d.Scheme
	u.Host = d.Host + ":" + d.Port
	u.Path = d.Index
	q := u.Query()
	q.Add("pretty", "true")
	u.RawQuery = q.Encode()
	fmt.Println("Request:", u)
	err := errors.New("")
	res, err := deleteHttpResource(u.String())
	if err != nil {
		return err.Error(), false
	}
	return res, true
}
