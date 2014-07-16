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

var Commands = struct {
	Version string
	Exit    string
	Help    string
	Host    string
	Port    string
	Index   string
	Dir     string
	Load    string
	Log     string
	Get     string
	Post    string
	Put     string
	Reindex string
}{
	Version: "version",
	Exit:    "exit",
	Help:    "help",
	Host:    "host",
	Port:    "port",
	Index:   "index",
	Dir:     "dir",
	Load:    "load",
	Log:     "log",
	Get:     "get",
	Post:    "post",
	Put:     "put",
	Reindex: "reindex",
}

var (
	server        = Server{}
	loadedRequest = LoadedRequest{}
	logLevel      = 0
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
	elRepl  ::  elasticsearch repl
	------------------------------
	`

	fmt.Println(message)
}

func reploop() {
	commandParser := NewCommandParser()
	for {
		fmt.Print("> ")
		in := bufio.NewReader(os.Stdin)
		entered, err := in.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			break
		}
		entry := strings.Trim(entered, "\t \r\n")
		if logLevel > 0 {
			log(entry, logLevel)
		}
		command, err := commandParser.Parse(entry)
		if err != nil {
			fmt.Println("Unable to parse command.")
		}
		output := dispatch(command)
		if len(output) > 0 {
			fmt.Println(output)
			fmt.Println("")
		}
	}
}

func dispatch(cmd *Command) string {
	if cmd == nil {
		return ""
	}

	switch cmd.Name {
	case Commands.Version:
		return handleVersion()
	case Commands.Help:
		return handleHelp()
	case Commands.Exit:
		return handleExit()
	case Commands.Host:
		if cmd.Args == "" {
			return handleServerGet()
		} else {
			return handleServerSet(cmd)
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
	case Commands.Load:
		return handleLoad(cmd)
	case Commands.Log:
		return handleLog(cmd)
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

func log(entry string, logLevel int) {
	f, err := os.OpenFile("elrepl.history.log", os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Unable to open file for logging:", err)
	}
	defer f.Close()
	f.WriteString(entry + "\r\n")
}
