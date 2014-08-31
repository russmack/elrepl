package main

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

type AliasCmd struct{}

func init() {
	h := NewHandler()
	h.CommandName = "alias"
	h.CommandPattern = "(alias)( )(.*)"
	h.Usage = "alias (create host port indexName aliasName) | (remove host port indexName aliasName) | (move host port fromIndex toIndex aliasName)"
	h.CommandParser = func(cmd *Command) (string, bool) {
		pattFn := map[*regexp.Regexp]func([]string) (string, bool){
			// List all aliases on host
			regexp.MustCompile(`^alias$`): func(s []string) (string, bool) {
				d := Resource{
					Endpoint: "_aliases",
					Scheme:   "http",
					Host:     server.host,
					Port:     server.port,
				}
				c := AliasCmd{}
				r, ok := c.GetAll(d)
				return r, ok
			},
			// Alias help
			regexp.MustCompile(`^alias /\?$`): func(s []string) (string, bool) {
				return "", false
			},
			// Create alias
			regexp.MustCompile(`^alias create ([a-zA-Z0-9\.\-]+) ([0-9]{1,5}) ([a-zA-Z0-9\.\-]+) ([a-zA-Z0-9\.\-]+)$`): func(s []string) (string, bool) {
				d := Resource{
					Endpoint: "_aliases",
					Index:    s[3],
					Alias:    s[4],
				}
				c := AliasCmd{}
				r, ok := c.Create(d)
				return r, ok
			},
			// Remove alias
			regexp.MustCompile(`^alias remove ([a-zA-Z0-9\.\-]+) ([0-9]{1,5}) ([a-zA-Z0-9\.\-]+) ([a-zA-Z0-9\.\-]+)$`): func(s []string) (string, bool) {
				d := Resource{
					Endpoint: "_aliases",
					Index:    s[3],
					Alias:    s[4],
				}
				c := AliasCmd{}
				r, ok := c.Remove(d)
				return r, ok
			},
			// Move alias
			regexp.MustCompile(`^alias move ([a-zA-Z0-9\.\-]+) ([0-9]{1,5}) ([a-zA-Z0-9\.\-]+) ([a-zA-Z0-9\.\-]+) ([a-zA-Z0-9\.\-]+)$`): func(s []string) (string, bool) {
				dFrom := Resource{
					Endpoint: "_aliases",
					Index:    s[3],
					Alias:    s[5],
				}
				dTarget := Resource{
					Index: s[4],
					Alias: s[5],
				}
				c := AliasCmd{}
				r, ok := c.Move(dFrom, dTarget)
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

func (c *AliasCmd) GetAll(d Resource) (string, bool) {
	if server.host == "" || server.port == "" {
		return "Missing host or port environment config.", false
	}
	u := new(url.URL)
	u.Scheme = d.Scheme
	u.Host = d.Host + ":" + d.Port
	u.Path = "/" + d.Endpoint
	q := u.Query()
	q.Add("pretty", "true")
	u.RawQuery = q.Encode()
	fmt.Println("Request:", u)
	err := errors.New("")
	res, err := getHttpResource(u.String())
	if err != nil {
		return err.Error(), false
	}
	return res, true
}

func (c *AliasCmd) Create(d Resource) (string, bool) {
	if server.host == "" || server.port == "" {
		return "Missing host or port environment config.", false
	}
	//curl -XPOST "http://10.1.1.12:9200/_aliases" -d '{ "actions": [ { "add": { "index": "podcasts-2014-07-29-001", "alias": "podcastsupdater" } } ] }'
	//post _alias?pretty { "actions": [ { "add": { "index": "podcasts-2014-07-29-001", "alias": "podcastsupdater" } } ] }
	urlString := "post " + d.Endpoint + " " + "{\"actions\": [ { \"add\": { \"index\": \"" + d.Index + "\", \"alias\": \"" + d.Alias + "\" } } ] }"
	cmdParser := NewCommandParser()
	newCmd, err := cmdParser.Parse(urlString)
	if err != nil {
		return err.Error(), false
	}
	dispatcher := NewDispatcher()
	res := dispatcher.Dispatch(newCmd)
	return res, true
}

func (c *AliasCmd) Remove(d Resource) (string, bool) {
	if server.host == "" || server.port == "" {
		return "Missing host or port environment config.", false
	}
	//curl -XPOST "http://10.1.1.12:9200/_aliases" -d '{ "actions": [ { "remove": { "index": "podcasts-2014-05-07-0103", "alias": "podcastsupdater" } } ] }'
	//post _alias?pretty { "actions": [ { "remove": { "index": "podcasts-2014-05-07-0103", "alias": "podcastsupdater" } } ] }
	urlString := "post " + d.Endpoint + " " + "{\"actions\": [ { \"remove\": { \"index\": \"" + d.Index + "\", \"alias\": \"" + d.Alias + "\" } } ] }"
	cmdParser := NewCommandParser()
	newCmd, err := cmdParser.Parse(urlString)
	if err != nil {
		return err.Error(), false
	}
	dispatcher := NewDispatcher()
	res := dispatcher.Dispatch(newCmd)
	return res, true
}

func (c *AliasCmd) Move(dFrom Resource, dTarget Resource) (string, bool) {
	if server.host == "" || server.port == "" {
		return "Missing host or port environment config.", false
	}
	postData := "{ \"actions\": [ { \"remove\": { \"alias\": \"" + dFrom.Alias + "\", \"index\": \"" + dFrom.Index + "\" }}, { \"add\": { \"alias\": \"" + dTarget.Alias + "\", \"index\": \"" + dTarget.Index + "\" } } ] }"
	urlString := "post " + dFrom.Endpoint + " " + postData
	cmdParser := NewCommandParser()
	newCmd, err := cmdParser.Parse(urlString)
	if err != nil {
		return err.Error(), false
	}
	dispatcher := NewDispatcher()
	res := dispatcher.Dispatch(newCmd)
	return res, true
}
