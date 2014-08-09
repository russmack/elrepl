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

var (
	server          = Server{}
	loadedRequest   = LoadedRequest{}
	logLevel        = 0
	HandlerRegistry = make(map[string]*Handler)
)

func init() {
	server.host = "localhost"
	server.port = "9200"
}

func main() {
	displayWelcome()
	//displayAvailableHandlers()
	reploop()
}

func displayWelcome() {
	message := `
	elRepl  ::  elasticsearch repl
	------------------------------
	`
	fmt.Println(message)
}

func displayAvailableHandlers() {
	fmt.Println("Handlers available:")
	for k, _ := range HandlerRegistry {
		fmt.Println(k)
	}
}

func reploop() {
	commandParser := NewCommandParser()
	dispatcher := NewDispatcher()
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
		output := dispatcher.Dispatch(command)
		if len(output) > 0 {
			fmt.Println(output)
			fmt.Println("")
		}
	}
}

func log(entry string, logLevel int) {
	f, err := os.OpenFile("elrepl.history.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Unable to open file for logging:", err)
	}
	defer f.Close()
	_, err = f.WriteString(entry + "\r\n")
	if err != nil {
		fmt.Println("Unable to write to log file:", err)
	}
}

func usageMessage(msg string) string {
	return "Usage: " + msg
}
