// elRepl is a repl for elasticsearch.
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Server struct {
	host  string
	port  string
	index string
}

type LoadedRequest struct {
	request string
}

const (
	CommandServer  = "host"
	CommandPort    = "port"
	CommandIndex   = "index"
	CommandDir     = "dir"
	CommandLoad    = "load"
	CommandGet     = "get"
	CommandPost    = "post"
	CommandPut     = "put"
	CommandReindex = "reindex"
)

var (
	server        = Server{}
	loadedRequest = LoadedRequest{}
)

func init() {
	server.host = "localhost"
	server.port = "9200"
}

func main() {
	displayWelcome()
	reploop()
}

func displayWelcome() {
	message := `
	el Repl
	=======

	Welcome to el Repl, an elasticsearch repl.
	`

	fmt.Println(message)
}

func reploop() {
	for {
		fmt.Print("> ")
		in := bufio.NewReader(os.Stdin)
		entered, err := in.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			break
		}
		entry := strings.Trim(entered, "\t \r\n")
		output := dispatch(entry)
		if len(output) > 0 {
			fmt.Println(output)
			fmt.Println("")
		}
	}
}

func dispatch(entry string) string {
	if len(entry) == 0 {
		return entry
	}
	switch {
	case entry == "?" || entry == "help":
		return handleHelp()
	case entry == "version":
		return handleVersion()
	case entry == "exit" || entry == "quit" || entry == "bye":
		return handleExit()
	case strings.HasPrefix(entry, CommandServer+" "):
		return handleServerSet(entry)
	case entry == CommandServer:
		return handleServerGet()
	case strings.HasPrefix(entry, CommandPort+" "):
		return handlePortSet(entry)
	case entry == CommandPort:
		return handlePortGet()
	case strings.HasPrefix(entry, CommandIndex+" "):
		return handleIndexSet(entry)
	case entry == CommandIndex:
		return handleIndexGet()
	case strings.HasPrefix(entry, CommandDir):
		return handleDir(entry)
	case strings.HasPrefix(entry, CommandLoad):
		return handleLoad(entry)
	case strings.HasPrefix(entry, CommandGet):
		return handleGet(entry)
	case strings.HasPrefix(entry, CommandPost):
		return handlePost(entry)
	case strings.HasPrefix(entry, CommandPut):
		return handlePut(entry)
	case strings.HasPrefix(entry, CommandReindex):
		return handleReindex(entry)
	default:
		return handleUnknownEntry(entry)
	}
}
