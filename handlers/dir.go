package handlers

import (
	"fmt"
	"github.com/russmack/elrepl/types"
	"io/ioutil"
)

func init() {
	fmt.Println("REGISTERING HANDLER: DIR")
	h := types.NewHandler()
	h.CommandName = "dir"
	h.CommandPattern = "(dir)( )(.*)"
	h.HandlerFunc = func(cmd *types.Command) string {
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
