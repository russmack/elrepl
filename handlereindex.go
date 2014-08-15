package main

import (
	"fmt"
	"github.com/mattbaird/elastigo/api"
	"github.com/mattbaird/elastigo/core"
)

// reindex localhost:9200/srcindex/type localhost:9200/targetindex/routing
func init() {
	h := NewHandler()
	h.CommandName = "reindex"
	h.CommandPattern = "(reindex)( )(.*)"
	h.Usage = "reindex sourceHost port sourceIndex sourceType targetHost port targetIndex [routing]"
	h.CommandParser = func(cmd *Command) (ParseMap, bool) {
		p := ParseMap{}
		argsCount := len(cmd.Args)
		if argsCount == 7 || argsCount == 8 {
			p["srcHost"] = cmd.Args[0]
			p["srcPort"] = cmd.Args[1]
			p["srcIndex"] = cmd.Args[2]
			p["srcType"] = cmd.Args[3]
			p["tgtHost"] = cmd.Args[4]
			p["tgtPort"] = cmd.Args[5]
			p["tgtIndex"] = cmd.Args[6]
			if argsCount == 8 {
				p["tgtRouting"] = cmd.Args[7]
			}
		} else {
			//case "/?"
			//case ""
			return p, false
		}
		return p, true
	}
	h.HandlerFunc = func(cmd *Command) string {
		fmt.Println("Reindexing...")
		p, ok := h.CommandParser(cmd)
		if !ok {
			return usageMessage(h.Usage)
		}
		srcHost := p["srcHost"]
		srcPort := p["srcPort"]
		srcIndex := p["srcIndex"]
		srcType := p["srcType"]
		tgtHost := p["tgtHost"]
		tgtPort := p["tgtPort"]
		tgtIndex := p["tgtIndex"]
		tgtRouting := p["tgtRouting"]

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
