package main

import (
//"fmt"
)

type Handler struct {
	CommandName    string
	CommandPattern string
	HandlerFunc    func(*Command) string
}

func NewHandler() *Handler {
	return &Handler{}
}
