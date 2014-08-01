package main

import (
	"github.com/russmack/elrepl/handlers"
	"github.com/russmack/elrepl/types"
)

type Dispatcher struct {
	//
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{}
}

func (d *Dispatcher) Dispatch(cmd *types.Command) string {
	if cmd == nil {
		return ""
	}

	s := handlers.HandlerRegistry[cmd.Name].HandlerFunc(cmd)

	return "cool: " + s

	switch cmd.Name {
	case Commands.Version:
		return handleVersion()
	case Commands.Help:
		return handleHelp()
	case Commands.Exit:
		return handleExit()
	case Commands.Host:
		if cmd.Args == "" {
			return handleHostGet()
		} else {
			return handleHostSet(cmd)
		}
	case Commands.Port:
		if cmd.Args == "" {
			return handlePortGet()
		} else {
			return handlePortSet(cmd)
		}
	case Commands.Index:
		if cmd.Args == "" {
			return handleIndexGet()
		} else {
			return handleIndexSet(cmd)
		}
	case Commands.Dir:
		return handleDir(cmd)
	case Commands.Log:
		return handleLog(cmd)
	case Commands.Load:
		return handleLoad(cmd)
	case Commands.Run:
		return handleRun(cmd)
	case Commands.Get:
		return handleGet(cmd)
	case Commands.Post:
		return handlePost(cmd)
	case Commands.Put:
		return handlePut(cmd)
	case Commands.Reindex:
		return handleReindex(cmd)
	default:
		return handleUnknownEntry(cmd)
	}
}
