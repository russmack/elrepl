package commands

import (
	"errors"
	"fmt"
	"github.com/russmack/elrepl/backend"
	"net/url"
	"regexp"
	"strings"
)

type AliasCmd struct{}

func init() {
	h := backend.NewHandler()
	h.CommandName = "alias"
	h.CommandPattern = "(alias)( )(.*)"
	h.Usage = "alias\n" +
		"  alias host port\n" +
		"  alias create host port indexName aliasName\n" +
		"  alias remove host port indexName aliasName\n" +
		"  alias move host port fromIndex toIndex aliasName"
	h.CommandParser = func(cmd *backend.Command) (string, bool) {
		pattFn := map[*regexp.Regexp]func([]string) (string, bool){
			// List all aliases on host
			regexp.MustCompile(`^alias$`): func(s []string) (string, bool) {
				d := backend.Resource{
					Endpoint: "_aliases",
					Scheme:   "http",
					Host:     backend.Server.Host,
					Port:     backend.Server.Port,
				}
				c := AliasCmd{}
				r, ok := c.GetAll(d)
				return r, ok
			},
			regexp.MustCompile(`^alias ([a-zA-Z0-9\.\-]+) ([0-9]{1,5})$`): func(s []string) (string, bool) {
				d := backend.Resource{
					Endpoint: "_aliases",
					Scheme:   "http",
					Host:     s[1],
					Port:     s[2],
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
				d := backend.Resource{
					Endpoint: "_aliases",
					Host:     s[1],
					Port:     s[2],
					Index:    s[3],
					Alias:    s[4],
				}
				c := AliasCmd{}
				r, ok := c.Create(d)
				return r, ok
			},
			// Remove alias
			regexp.MustCompile(`^alias remove ([a-zA-Z0-9\.\-]+) ([0-9]{1,5}) ([a-zA-Z0-9\.\-]+) ([a-zA-Z0-9\.\-]+)$`): func(s []string) (string, bool) {
				d := backend.Resource{
					Endpoint: "_aliases",
					Host:     s[1],
					Port:     s[2],
					Index:    s[3],
					Alias:    s[4],
				}
				c := AliasCmd{}
				r, ok := c.Remove(d)
				return r, ok
			},
			// Move alias
			regexp.MustCompile(`^alias move ([a-zA-Z0-9\.\-]+) ([0-9]{1,5}) ([a-zA-Z0-9\.\-]+) ([a-zA-Z0-9\.\-]+) ([a-zA-Z0-9\.\-]+)$`): func(s []string) (string, bool) {
				dFrom := backend.Resource{
					Endpoint: "_aliases",
					Host:     s[1],
					Port:     s[2],
					Index:    s[3],
					Alias:    s[5],
				}
				dTarget := backend.Resource{
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

func (c *AliasCmd) GetAll(d backend.Resource) (string, bool) {
	if d.Host == "" || d.Port == "" {
		return "Missing host or port.", false
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
	res, err := backend.GetHttpResource(u.String())
	if err != nil {
		return err.Error(), false
	}
	return res, true
}

func (c *AliasCmd) Create(d backend.Resource) (string, bool) {
	if d.Host == "" || d.Port == "" {
		return "Missing host or port.", false
	}
	postData := "{\"actions\": [ { \"add\": { \"index\": \"" + d.Index + "\", \"alias\": \"" + d.Alias + "\" } } ] }"
	urlString := "post " + d.Host + " " + d.Port + " " + d.Endpoint + " " + postData
	cmdParser := backend.NewCommandParser()
	newCmd, err := cmdParser.Parse(urlString)
	if err != nil {
		return err.Error(), false
	}
	dispatcher := backend.NewDispatcher()
	res, ok := dispatcher.Dispatch(newCmd)
	return res, ok
}

func (c *AliasCmd) Remove(d backend.Resource) (string, bool) {
	if d.Host == "" || d.Port == "" {
		return "Missing host or port.", false
	}
	postData := "{\"actions\": [ { \"remove\": { \"index\": \"" + d.Index + "\", \"alias\": \"" + d.Alias + "\" } } ] }"
	urlString := "post " + d.Host + " " + d.Port + " " + d.Endpoint + " " + postData
	cmdParser := backend.NewCommandParser()
	newCmd, err := cmdParser.Parse(urlString)
	if err != nil {
		return err.Error(), false
	}
	dispatcher := backend.NewDispatcher()
	res, ok := dispatcher.Dispatch(newCmd)
	return res, ok
}

func (c *AliasCmd) Move(dFrom backend.Resource, dTarget backend.Resource) (string, bool) {
	if dFrom.Host == "" || dFrom.Port == "" {
		return "Missing host or port.", false
	}
	postData := "{ \"actions\": [ { \"remove\": { \"alias\": \"" + dFrom.Alias + "\", \"index\": \"" + dFrom.Index + "\" }}, { \"add\": { \"alias\": \"" + dTarget.Alias + "\", \"index\": \"" + dTarget.Index + "\" } } ] }"
	urlString := "post " + dFrom.Host + " " + dFrom.Port + " " + dFrom.Endpoint + " " + postData
	cmdParser := backend.NewCommandParser()
	newCmd, err := cmdParser.Parse(urlString)
	if err != nil {
		return err.Error(), false
	}
	dispatcher := backend.NewDispatcher()
	res, ok := dispatcher.Dispatch(newCmd)
	return res, ok
}
