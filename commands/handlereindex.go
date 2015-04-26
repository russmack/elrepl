package commands

import (
	"fmt"
	"github.com/mattbaird/elastigo/lib"
	"github.com/russmack/elrepl/backend"
	"regexp"
	"strings"
)

type ReindexCmd struct{}

// reindex localhost:9200/srcindex/type localhost:9200/targetindex/routing
func init() {
	h := backend.NewHandler()
	h.CommandName = "reindex"
	h.CommandPattern = "(reindex)( )(.*)"
	h.Usage = "reindex sourceHost sourcePort sourceIndex sourceType targetHost targetPort targetIndex [routing]"
	h.CommandParser = func(cmd *backend.Command) (string, bool) {
		pattFn := map[*regexp.Regexp]func([]string) (string, bool){
			// Index help
			regexp.MustCompile(`^reindex /\?$`): func(s []string) (string, bool) {
				return "", false
			},
			// Reindex
			regexp.MustCompile(`^reindex ([a-zA-Z0-9\.\-]+) ([0-9]{1,5}) ([a-zA-Z0-9\.\-]+) ([a-zA-Z0-9\.\-]+) ([a-zA-Z0-9\.\-]+) ([0-9]{1,5}) ([a-zA-Z0-9\.\-]+) ([a-zA-Z0-9\.\-]+)?$`): func(s []string) (string, bool) {
				dSource := backend.Resource{
					Scheme: "http",
					Host:   s[1],
					Port:   s[2],
					Index:  s[3],
					Type:   s[4],
				}
				dTarget := backend.Resource{
					Scheme:  "http",
					Host:    s[5],
					Port:    s[6],
					Index:   s[7],
					Routing: s[8],
				}
				c := ReindexCmd{}
				r, ok := c.Reindex(dSource, dTarget)
				return r, ok
			},
		}
		r, ok := h.Tokenize(strings.TrimSpace(cmd.Instruction), pattFn)
		return r, ok
	}
	h.HandlerFunc = func(cmd *backend.Command) (string, bool) {
		fmt.Println("Reindexing...")
		r, ok := h.CommandParser(cmd)
		if !ok {
			if r != "" {
				r += "\n\n"
			}
			return h.Usage, false
		}
		return r, true
	}
	backend.HandlerRegistry[h.CommandName] = h
}

func (c *ReindexCmd) Reindex(dSource backend.Resource, dTarget backend.Resource) (string, bool) {
	api := elastigo.NewConn()
	api.Domain = dSource.Host
	api.Port = dSource.Port

	fmt.Println("Scanning...")
	scanArgs := map[string]interface{}{"search_type": "scan", "scroll": "1m", "size": "1000"}
	scanResult, err := api.SearchUri(dSource.Index, dSource.Type, scanArgs)
	if err != nil {
		fmt.Println("Failed getting scan result for index:", dSource.Index, "; err:", err)
		return err.Error(), false
	}

	//total := scanResult.Hits.Total

	scrollId := scanResult.ScrollId
	counter := 0
	failures := 0

	fmt.Println("Scrolling...")
	scrollArgs := map[string]interface{}{"scroll": "1m"}
	scrollResult, err := api.Scroll(scrollArgs, scrollId)
	if err != nil {
		fmt.Println("Failed getting scroll result for index:", dSource.Index, "; err:", err)
		return err.Error(), false
	}

	fmt.Println("Indexing...")
	var indexArgs map[string]interface{} = nil
	if dTarget.Routing != "" {
		indexArgs = map[string]interface{}{"routing": dTarget.Routing}
	}
	for len(scrollResult.Hits.Hits) > 0 {
		fmt.Println("Scroll result hits:", len(scrollResult.Hits.Hits))
		for _, j := range scrollResult.Hits.Hits {
			api.Domain = dTarget.Host
			api.Port = dTarget.Port

			_, err := api.Index(dTarget.Index, dSource.Type, j.Id, indexArgs, j.Source)
			if err != nil {
				fmt.Println("Failed inserting document, id:", j.Id, "; ", err)
				failures++
				continue
			}
			counter++
		}

		api.Domain = dSource.Host
		api.Port = dSource.Port
		// ScrollId changes with every request.
		scrollId = scrollResult.ScrollId
		scrollArgs := map[string]interface{}{"scroll": "1m"}
		scrollResult, err = api.Scroll(scrollArgs, scrollId)
		if err != nil {
			fmt.Println("Failed getting scroll result for index:", dSource.Index, "; err:", err)
			return err.Error(), false
		}
	}
	return fmt.Sprintf("Total processed: %d.  %d failed.", counter, failures), true
}
