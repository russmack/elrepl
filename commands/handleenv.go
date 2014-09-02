package commands

import (
	"github.com/russmack/elrepl/backend"
	"regexp"
	"strings"
)

type EnvCmd struct{}

func init() {
	h := backend.NewHandler()
	h.CommandName = "env"
	h.CommandPattern = "(env)(( )(.*))"
	h.Usage = "env [ (host hostname) | (port portnumber) | (index indexname ) ]"
	h.CommandParser = func(cmd *backend.Command) (string, bool) {
		pattFn := map[*regexp.Regexp]func([]string) (string, bool){
			// Get env
			regexp.MustCompile(`^env$`): func(s []string) (string, bool) {
				d := backend.Resource{}
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
				d := backend.Resource{
					Host: s[1],
				}
				c := EnvCmd{}
				r, ok := c.SetHost(d)
				return r, ok
			},
			// Set env port
			regexp.MustCompile(`^env port ([0-9]{1,5})$`): func(s []string) (string, bool) {
				d := backend.Resource{
					Port: s[1],
				}
				c := EnvCmd{}
				r, ok := c.SetPort(d)
				return r, ok
			},
			// Set env index
			regexp.MustCompile(`^env index ([a-zA-Z0-9\.\-]+)$`): func(s []string) (string, bool) {
				d := backend.Resource{
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

func (c *EnvCmd) Get(d backend.Resource) (string, bool) {
	return `
	elRepl version 0.1

	Host: ` + backend.Server.Host + `
	Port: ` + backend.Server.Port + `
	Index: ` + backend.Server.Index + `
	`, true
}

func (c *EnvCmd) SetHost(d backend.Resource) (string, bool) {
	backend.Server.Host = d.Host
	return "Host set: " + backend.Server.Host, true
}

func (c *EnvCmd) SetPort(d backend.Resource) (string, bool) {
	backend.Server.Port = d.Port
	return "Port set: " + backend.Server.Port, true
}

func (c *EnvCmd) SetIndex(d backend.Resource) (string, bool) {
	backend.Server.Index = d.Index
	return "Index set: " + backend.Server.Index, true
}
