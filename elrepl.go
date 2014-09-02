// elRepl is a repl for elasticsearch.
package main

import (
	"bufio"
	"fmt"
	"github.com/russmack/elrepl/backend"
	_ "github.com/russmack/elrepl/commands"
	"os"
	"strings"
)

func init() {
	backend.Server.Host = "localhost"
	backend.Server.Port = "9200"
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
	for k, _ := range backend.HandlerRegistry {
		fmt.Println(k)
	}
}

func reploop() {
	commandParser := backend.NewCommandParser()
	dispatcher := backend.NewDispatcher()
	for {
		fmt.Print("> ")
		in := bufio.NewReader(os.Stdin)
		entered, err := in.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			break
		}
		entry := strings.Trim(entered, "\t \r\n")
		if backend.LogLevel > 0 {
			log(entry, backend.LogLevel)
		}
		command, err := commandParser.Parse(entry)
		if err != nil {
			fmt.Println("Unable to parse command.")
		}
		output, ok := dispatcher.Dispatch(command)
		if len(output) > 0 {
			if !ok {
				fmt.Println(usageMessage(output))
			} else {
				fmt.Println(output)
			}
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
