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

type ESRequest struct {
	Size int `json:"size"`
	Sort struct {
		Timestamp struct {
			Order        string `json:"order"`
			UnmappedType string `json:"unmapped_type"`
		} `json:"@timestamp"`
	}
	Query struct {
		Filtered struct {
			Query struct {
				QueryString struct {
					AnalyzeWildcard string `json:"analyze_wildcard"`
					Query           string `json:"query"`
				} `json:"query_string"`
			} `json:"query"`
		} `json:"filered"`
		Filter struct {
			Bool struct {
				Must    []ESMust      `json:"must"`
				MustNot []interface{} `json:"must_not"`
			} `json:"bool"`
		} `json:"filter"`
	} `json:"query"`
}

type ESMust struct {
	Range struct {
		Timestamp struct {
			Gte interface{} `json:"gte"`
			Lte interface{} `json:"lte"`
		} `json:"@timestamp"`
	} `json:"range"`
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
	var (
		gte Gte
		lte = time.Now().Add(time.Duration(-c.StartTime) * time.Minute)
	)

	for syncCount := 0; syncCount >= 0; syncCount++ {
		// Set time: last 10min or last sync_interval
		if syncCount > 0 {
			gte.Time = lte.Add(time.Duration(-c.SyncInterval) * time.Second)
		} else {
			gte.Time = lte.Add(time.Duration(-c.SyncDepth) * time.Minute)
		}

		var (
			esrequest = ESRequest{}
			must      = ESMust{}
			response  = Es_resp{}
		)

		must.Range.Timestamp.Gte = gte.Time
		must.Range.Timestamp.Lte = lte

		esrequest.Sort.Timestamp.Order = "asc"
		esrequest.Sort.Timestamp.UnmappedType = "long"
		esrequest.Query.Filtered.Query.QueryString.AnalyzeWildcard = "true"
		esrequest.Query.Filtered.Query.QueryString.Query = string(service)

		esrequest.Query.Filter.Bool.Must = []ESMust{must}

		jsonpost, err := json.MarshalIndent(&esrequest, "", "\t")
		errorhandler.BasicCheckOrExit(err)
		log.Debugf("Elastic Search Request:\n %s", string(jsonpost))

		// Craft the request URI
		uri_ary := []string{"http://", c.ElasticsearchURL, ":", c.ElasticsearchPort, "/_search?pretty"} //c.Elasticsearch_index, "/_search?pretty"}
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
							if c.SearchHost {
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
		log.Debug("Sync ", time.Duration(c.SyncInterval))
		time.Sleep(time.Second * time.Duration(c.SyncInterval))
	}
}

func SetLogger(verbose bool) {
	if verbose {
		log.SetLevel(log.DebugLevel)
		log.Debug("DEBUG Logger")
	} else {
		log.SetLevel(log.InfoLevel)
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
	SetLogger(config.LogVerbose)
	log.Debugf("%+v", config)
	query, err := config.GetDefinedQuery()
	errorhandler.BasicCheckOrExit(err)
	log.Infof("Starting new search for %s", query)
	// Roll into the query loop
	searchRunner(query, config)
}
