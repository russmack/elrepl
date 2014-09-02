package backend

import (
	"fmt"
)

type Dispatcher struct{}

type ServerEnv struct {
	Host  string
	Port  string
	Index string
}

type LoadedRequestEnv struct {
	Request string
}

var HandlerRegistry = make(map[string]*Handler)

var (
	Server        = ServerEnv{}
	LoadedRequest = LoadedRequestEnv{}
	LogLevel      = 0
)

func NewDispatcher() *Dispatcher {
	return &Dispatcher{}
}

func (d *Dispatcher) Dispatch(cmd *Command) (string, bool) {
	if cmd == nil {
		return "", false
	}

	h, ok := HandlerRegistry[cmd.Name]
	if !ok {
		return fmt.Sprintf("Command not found: %s", cmd.Name), ok
	} else {
		return h.HandlerFunc(cmd)
	}
}
