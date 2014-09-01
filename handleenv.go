package main

import (
	"regexp"
	"strings"
)

type EnvCmd struct{}

func init() {
	h := NewHandler()
	h.CommandName = "env"
	h.CommandPattern = "(env)(( )(.*))"
	h.Usage = "env [ (host hostname) | (port portnumber) | (index indexname ) ]"
	h.CommandParser = func(cmd *Command) (string, bool) {
		pattFn := map[*regexp.Regexp]func([]string) (string, bool){
			// Get env
			regexp.MustCompile(`^env$`): func(s []string) (string, bool) {
				d := Resource{}
				c := EnvCmd{}
				r, ok := c.Get(d)
				return r, ok
			},
			// Env help
			regexp.MustCompile(`^env /\?$`): func(s []string) (string, bool) {
				return "", false
			},
			// Set env host
			regexp.MustCompile(`^env host ([a-zA-Z0-9\.\-]+)$`): func(s []string) (string, bool) {
				d := Resource{
					Host: s[1],
				}
				c := EnvCmd{}
				r, ok := c.SetHost(d)
				return r, ok
			},
			// Set env port
			regexp.MustCompile(`^env port ([0-9]{1,5})$`): func(s []string) (string, bool) {
				d := Resource{
					Port: s[1],
				}
				c := EnvCmd{}
				r, ok := c.SetPort(d)
				return r, ok
			},
			// Set env index
			regexp.MustCompile(`^env index ([a-zA-Z0-9\.\-]+)$`): func(s []string) (string, bool) {
				d := Resource{
					Index: s[1],
				}
				c := EnvCmd{}
				r, ok := c.SetIndex(d)
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

func (c *EnvCmd) Get(d Resource) (string, bool) {
	return `
	elRepl version 0.1

	Host: ` + server.host + `
	Port: ` + server.port + `
	Index: ` + server.index + `
	`, true
}

func (c *EnvCmd) SetHost(d Resource) (string, bool) {
	server.host = d.Host
	return "Host set: " + server.host, true
}

func (c *EnvCmd) SetPort(d Resource) (string, bool) {
	server.port = d.Port
	return "Port set: " + server.port, true
}

func (c *EnvCmd) SetIndex(d Resource) (string, bool) {
	server.index = d.Index
	return "Index set: " + server.index, true
}
