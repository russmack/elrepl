package main

import (
	"fmt"
)

type Dispatcher struct {
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{}
}

func (d *Dispatcher) Dispatch(cmd *Command) string {
	if cmd == nil {
		return ""
	}

	h, ok := HandlerRegistry[cmd.Name]
	if !ok {
		return fmt.Sprintf("Command not found: %s", cmd.Name)
	} else {
		return h.HandlerFunc(cmd)
	}
}
