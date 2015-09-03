package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/mgutz/ansi"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

// Config struct
type Config struct {
	Define              map[string]string `yaml:"define"`
	Sync_interval       int               `yaml:"sync_interval"`
	Sync_depth          int               `yaml:"sync_depth"`
	Elasticsearch_url   string            `yaml:"elasticsearch_url"`
	Elasticsearch_port  string            `yaml:"elasticsearch_port"`
	Elasticsearch_index string            `yaml:"elasticsearch_index"`
	Host                bool
	Highlight           bool
	StartTime           time.Time
	Count               int
}

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

// Define flag overrides
var home = os.Getenv("HOME")
var config_path = flag.String("c", strings.Join([]string{home, "/.loga/config.yaml"}, ""), "The path to the config.yaml.")
var define_service = flag.String("d", "", "A one-time defined service. Must be valid ES query.")
var elastic_url = flag.String("e", "", "Elastic search URL.")
var sync_interval = flag.Int("si", 5, "Query interval in seconds.")
var sync_depth = flag.Int("sd", 10, "Sync Depth - how far back to go on initial query.")
var elastic_port = flag.String("p", "", "Elastic Search port.")
var elastic_index = flag.String("in", "", "Elastic Search index.")
var verbose = flag.Bool("v", false, "Verbosity.")
var service = flag.String("s", "", "Query already defined service in config.yaml.")
var srch_host = flag.Bool("h", false, "Specific hostname to search.")
var highlight = flag.Bool("hl", false, "Highlight the string with the query.")
var startTime = flag.Int("st", 0, "Start time for query in minutes. Ex: -st 20 starts query 20 minutes ago.")
var count = flag.Int("co", 500, "The number of results to return.")

func options(config_path string) (o Config, err error) {
	config_file, err := ioutil.ReadFile(config_path)
	if err != nil {
		log.Error("Are you sure ~/.loga/config.yaml exists?")
		panic(err)
	}
	// Unmarshal the config to a map
	err = yaml.Unmarshal(config_file, &o)
	if err != nil {
		panic(err)
	}
	// Override a million things
	if len(*elastic_url) > 1 {
		o.Elasticsearch_url = *elastic_url
	}
	// Sync interval in seconds
	o.Sync_interval = *sync_interval
	// Port on ES to use
	if len(*elastic_port) > 0 {
		o.Elasticsearch_port = *elastic_port
	}
	if len(*elastic_index) > 0 {
		o.Elasticsearch_index = *elastic_index
	}
	// Sync depth to return
	o.Sync_depth = *sync_depth
	// Set host in line
	if *srch_host {
		o.Host = *srch_host
	}
	// Highlight query in output
	if *highlight {
		o.Highlight = *highlight
	}
	// Configure start time for query
	now := time.Now()
	if *startTime > 0 {
		o.StartTime = now.Add(time.Duration(-*startTime) * time.Minute)
	} else {
		o.StartTime = now
	}
	// Count of results to return
	o.Count = *count
	return o, nil
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

func query(service string, c Config) {
	for syncCount := 0; syncCount >= 0; syncCount++ {
		var gte Gte
		// Set time: last 10min or last sync_interval
		lte := time.Now() //c.StartTime
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

func lookup(defines map[string]string) (query string) {
	// If the -d flag is passed, test and return, if not, lookup, if not found, return err
	if len(*define_service) > 0 {
		log.Debug("Query for defined service: ", *define_service)
		//for i,v := range strings.Split(*define_service, "") {
		query := *define_service //fmt.Sprintf("\\%s\\", *define_service)
		log.Info(query)
		return query
	} else if len(*service) > 0 {
		if _, ok := defines[*service]; ok {
			log.Debug("Query for service", *service, "found in config:", defines[*service])
			query := defines[*service]
			return query
		} else {
			log.Error("Service", *service, "not found in config.")
			return "Service not found"
		}
	}
	return "nil"
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
	// Get cli flags
	flag.Parse()
	// Set loglevel
	if *verbose {
		log.SetLevel(log.DebugLevel)
		log.Debug("Loglevel: Debug")
	} else {
		log.SetLevel(log.InfoLevel)
		log.Info("Loglevel: Info")
	}
	// Set configuration, override config.yaml with flags
	config, err := options(*config_path)
	if err != nil {
		log.Error(err)
	}
	// Debug some things
	log.Debug("Configuration:")
	log.Debug(
		"Defines: ", config.Define, "\n",
		"Sync: ", config.Sync_interval, "\n",
		"ES URL: ", config.Elasticsearch_url, "\n",
		"ES Port: ", config.Elasticsearch_port, "\n",
		"ES Index: ", config.Elasticsearch_index, "\n",
		"Host ", config.Host, "\n",
		"Start Time ", config.StartTime, "\n",
		"Count ", config.Count, "\n",
		"Highlight ", config.Highlight)
	// Make sure the define set on the CLI exist if neccessary
	defines := config.Define
	svc_query := lookup(defines)
	log.Info("Querying ", *service, ": ", svc_query)
	// Roll into the query loop
	query(svc_query, config)
}
