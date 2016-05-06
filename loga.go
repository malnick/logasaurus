package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/malnick/logasaurus/config"
	"github.com/malnick/logasaurus/errorhandler"
	"github.com/mgutz/ansi"
)

type Es_resp struct {
	Hits interface{}
}

type Es_post struct {
	Size  int                          `json:"size"`
	Sort  map[string]map[string]string `json:"sort"`
	Query map[string]interface{}       `json:"query"`
}

type Gte struct {
	Time time.Time
}

func highlightQuery(line string, query string) {
	// Split query into multiple parts for regex
	q := strings.Split(query, " ")
	// Match the string
	match, err := regexp.Compile(q[0])
	if err != nil {
		panic(err)
	}

	// Split our line into an ary
	lineAry := strings.Split(line, " ")
	// Iterate the ary, finding the string match
	for i, s := range lineAry {
		if match.MatchString(s) {
			// Color just the string which matches
			hlQuery := ansi.Color(s, "yellow:black")
			// Thren break down into three parts
			lpt1 := lineAry[:i]
			lpt2 := lineAry[i:]
			lpt2 = append(lpt2[:0], lpt2[1:]...)
			// Contatenate back together
			part1 := strings.Join(lpt1, " ")
			part2 := strings.Join(lpt2, " ")
			final := []string{part1, hlQuery, part2}
			finalHl := strings.Join(final, " ")
			// Print the final output
			//log.Info(finalHl)
			fmt.Println(finalHl)
		}
	}
}

func searchRunner(service string, c config.Config) {
	for syncCount := 0; syncCount >= 0; syncCount++ {
		var gte Gte
		// Set time: last 10min or last sync_interval
		lte := c.StartTime
		if syncCount > 0 {
			gte.Time = lte.Add(time.Duration(-c.Sync_interval) * time.Second)
		} else {
			gte.Time = lte.Add(time.Duration(-c.Sync_depth) * time.Minute)
		}
		// Elasticsearch response
		var response Es_resp
		// The JSON query
		sort := map[string]map[string]string{
			"@timestamp": map[string]string{
				"order":         "asc",
				"unmapped_type": "long",
			},
		}
		query := map[string]interface{}{
			"filtered": map[string]interface{}{
				"query": map[string]map[string]interface{}{
					"query_string": {
						"query": string(service),
						//			"fields":           []string{"message", "host"},
						"analyze_wildcard": string("true"),
					},
				},
				"filter": map[string]map[string][]map[string]map[string]map[string]interface{}{
					"bool": {
						"must": {
							{
								"range": {
									"@timestamp": {
										"gte": gte.Time,
										"lte": lte,
									},
								},
							},
						},
						"must_not": {},
					},
				},
			},
		}

		postthis := Es_post{
			Size:  c.Count,
			Sort:  sort,
			Query: query,
		}
		log.Debug("ES Post Struc: ", postthis)

		jsonpost, err := json.Marshal(postthis)
		if err != nil {
			log.Error(err)
		}
		log.Debug("ES JSON Post: ", string(jsonpost))

		// Craft the request URI
		uri_ary := []string{"http://", c.Elasticsearch_url, ":", c.Elasticsearch_port, "/_search?pretty"} //c.Elasticsearch_index, "/_search?pretty"}
		query_uri := strings.Join(uri_ary, "")
		log.Debug("Query URI: ", query_uri)
		// Make request
		req, err := http.NewRequest("POST", query_uri, bytes.NewBuffer(jsonpost))
		if err != nil {
			log.Error(err)
			panic(err)
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		jsonRespBody, _ := ioutil.ReadAll(resp.Body)
		log.Debug("ES Response:")
		log.Debug(string(jsonRespBody))
		// Unmarshel json resp
		err = json.Unmarshal(jsonRespBody, &response)
		if err != nil {
			log.Error(err)
			panic(err)
		}
		// Print
		for k0, v0 := range response.Hits.(map[string]interface{}) {
			if k0 == "hits" {
				for _, v1 := range v0.([]interface{}) {
					for k2, v2 := range v1.(map[string]interface{}) {
						if k2 == "_source" {
							if c.Host {
								message := v2.(map[string]interface{})["message"].(string)
								host := ansi.Color(v2.(map[string]interface{})["host"].(string), "cyan:black")
								withHost := strings.Join([]string{host, " ", message}, "")
								if c.Highlight {
									highlightQuery(withHost, service)
								} else {
									//log.Info(logthis)
									fmt.Println(withHost)
								}
							} else {
								message := v2.(map[string]interface{})["message"].(string)
								if c.Highlight {
									highlightQuery(message, service)
								} else {
									//log.Info(message)
									fmt.Println(message)
								}
							}
						}
					}
				}
			}
		}
		log.Debug("Sync ", time.Duration(c.Sync_interval))
		time.Sleep(time.Second * time.Duration(c.Sync_interval))
	}
}

func main() {
	fmt.Println(`                        .       .                             `)
	fmt.Println(`                       / '.   .' \                            `)
	fmt.Println(`               .---.  <    > <    >  .---.                    `)
	fmt.Println(`               |    \  \ - ~ ~ - /  /    |                    `)
	fmt.Println(`               ~-..-~             ~-..-~                     `)
	fmt.Println(`            \~~~\.'                    './~~~/                `)
	fmt.Println(`  .-~~^-.    \__/                        \__/                 `)
	fmt.Println(`.'  O    \     /               /       \  \                   `)
	fmt.Println(`(_____'    \._.'              |         }  \/~~~/             `)
	fmt.Println(`  ----.         /       }     |        /    \__/              `)
	fmt.Println(`      \-.      |       /      |       /      \.,~~|           `)
	fmt.Println(`          ~-.__|      /_ - ~ ^|      /- _     \..-'   f: f:   `)
	fmt.Println(`               |     /        |     /     ~-.     -. _||_||_  `)
	fmt.Println(`               |_____|        |_____|         ~ - . _ _ _ _ _>`)
	fmt.Println(`██╗      ██████╗  ██████╗  █████╗ ███████╗ █████╗ ██╗   ██╗██████╗ ██╗   ██╗███████╗`)
	fmt.Println(`██║     ██╔═══██╗██╔════╝ ██╔══██╗██╔════╝██╔══██╗██║   ██║██╔══██╗██║   ██║██╔════╝`)
	fmt.Println(`██║     ██║   ██║██║  ███╗███████║███████╗███████║██║   ██║██████╔╝██║   ██║███████╗`)
	fmt.Println(`██║     ██║   ██║██║   ██║██╔══██║╚════██║██╔══██║██║   ██║██╔══██╗██║   ██║╚════██║`)
	fmt.Println(`███████╗╚██████╔╝╚██████╔╝██║  ██║███████║██║  ██║╚██████╔╝██║  ██║╚██████╔╝███████║`)
	fmt.Println(`╚══════╝ ╚═════╝  ╚═════╝ ╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝ ╚═════╝ ╚═╝  ╚═╝ ╚═════╝ ╚══════╝`)
	fmt.Println()
	config := config.ParseArgsReturnConfig()
	query, err := config.GetDefinedQuery()
	errorhandler.LogErrorAndExit(err)
	log.Infof("Starting new search for %s", query)
	// Roll into the query loop
	searchRunner(query, config)
}
