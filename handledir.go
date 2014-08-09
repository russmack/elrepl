package main

import (
	"io/ioutil"
)

func init() {
	h := NewHandler()
	h.CommandName = "dir"
	h.CommandPattern = "(dir)( )(.*)"
	h.Usage = "dir"
	h.HandlerFunc = func(cmd *Command) string {
		arg := cmd.Args
		if arg == "" {
			arg = "."
		}
		dirFiles, err := ioutil.ReadDir(arg)
		if err != nil {
			return err.Error()
		}
		files := ""
		for _, j := range dirFiles {
			files += j.Name() + "\n"
		}
		return files
	}
	HandlerRegistry[h.CommandName] = h
}
