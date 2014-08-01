package main

import (
	"fmt"
	"github.com/mattbaird/elastigo/api"
	"github.com/mattbaird/elastigo/core"
	"regexp"
)

// reindex localhost:9200/srcindex/type localhost:9200/targetindex/routing
func init() {
	h := NewHandler()
	h.CommandName = "reindex"
	h.CommandPattern = "(reindex)( )(.*)"
	h.HandlerFunc = func(cmd *Command) string {
		fmt.Println("Reindexing...")
		//args := strings.TrimPrefix(entry, CommandReindex+" ")
		args := cmd.Args

		// \w+|"[\w\s]*"
		r, err := regexp.Compile(`^(.*?):(\d+?)/(.*?)/(.*?)/? (.*?):(\d+?)/(.*?)(/(.*?))?$`)
		if err != nil {
			return err.Error()
		}
		fmt.Println("Parsing command...")
		matches := r.FindAllStringSubmatch(args, -1)[0]
		fmt.Println("Parsed matches:", len(matches))
		srcHost := matches[1]
		srcPort := matches[2]
		srcIndex := matches[3]
		srcType := matches[4]
		tgtHost := matches[5]
		tgtPort := matches[6]
		tgtIndex := matches[7]
		tgtRouting := matches[8]

		api.Domain = srcHost
		api.Port = srcPort

		fmt.Println("Scanning...")
		scanArgs := map[string]interface{}{"search_type": "scan", "scroll": "1m", "size": "1000"}
		scanResult, err := core.SearchUri(srcIndex, srcType, scanArgs)
		if err != nil {
			fmt.Println("Failed getting scan result for index:", srcIndex, "; err:", err)
			return err.Error()
		}

		//total := scanResult.Hits.Total

		scrollId := scanResult.ScrollId
		counter := 0
		failures := 0

		fmt.Println("Scrolling...")
		scrollArgs := map[string]interface{}{"scroll": "1m"}
		scrollResult, err := core.Scroll(scrollArgs, scrollId)
		if err != nil {
			fmt.Println("Failed getting scroll result for index:", srcIndex, "; err:", err)
			return err.Error()
		}

		fmt.Println("Indexing...")
		var indexArgs map[string]interface{} = nil
		if tgtRouting != "" {
			indexArgs = map[string]interface{}{"routing": tgtRouting}
		}
		for len(scrollResult.Hits.Hits) > 0 {
			fmt.Println("Scroll result hits:", len(scrollResult.Hits.Hits))
			for _, j := range scrollResult.Hits.Hits {
				api.Domain = tgtHost
				api.Port = tgtPort

				_, err := core.Index(tgtIndex, srcType, j.Id, indexArgs, j.Source)
				if err != nil {
					fmt.Println("Failed inserting document, id:", j.Id, "; ", err)
					failures++
					continue
				}
				counter++
			}

			api.Domain = srcHost
			api.Port = srcPort
			// ScrollId changes with every request.
			scrollId = scrollResult.ScrollId
			scrollArgs := map[string]interface{}{"scroll": "1m"}
			scrollResult, err = core.Scroll(scrollArgs, scrollId)
			if err != nil {
				fmt.Println("Failed getting scroll result for index:", srcIndex, "; err:", err)
				return err.Error()
			}
		}
		return fmt.Sprintf("Total processed: %d.  %d failed.", counter, failures)
	}
	HandlerRegistry[h.CommandName] = h
}
