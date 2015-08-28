package main

import (
	//	"encoding/json"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
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
var config_path = flag.String("c", "./config.yaml", "The path to the config.yaml.")
var define_service = flag.String("d", "", "A one-time defined service. Must be valid ES query.")
var elastic_url = flag.String("e", "", "Elastic search URL.")
var sync_interval = flag.Int("si", 0, "Query interval in seconds.")
var sync_depth = flag.Int("sd", 0, "Sync Depth - how far back to go on initial query.")
var elastic_port = flag.String("p", "", "Elastic Search port.")
var elastic_index = flag.String("in", "", "Elastic Search index.")
var verbose = flag.Bool("v", false, "Verbosity.")
var service = flag.String("s", "", "Query already defined service in config.yaml.")

func options(config_path string) (o Config, err error) {
	config_file, err := ioutil.ReadFile(config_path)
	if err != nil {
		return o, err
	}
	// Unmarshal the config to a map
	err = yaml.Unmarshal(config_file, &o)
	if err != nil {
		return o, err
	}
	// Override a million things
	if len(*elastic_url) > 1 {
		o.Elasticsearch_url = *elastic_url
	}
	if *sync_interval > 0 {
		o.Sync_interval = *sync_interval
	}
	if len(*elastic_port) > 0 {
		o.Elasticsearch_port = *elastic_port
	}
	if len(*elastic_index) > 0 {
		o.Elasticsearch_index = *elastic_index
	}
	if *sync_depth > 0 {
		o.Sync_depth = *sync_depth
	}
	return o, nil
}

func query(service string, c Config) {
	for syncCount := 0; syncCount >= 0; syncCount++ {
		var gte Gte
		// Set GTE time: last 10min or last sync_interval
		lte := time.Now()
		if syncCount > 0 {
			log.Debug("SYNC COUNT gt o")
			gte.Time = lte.Add(time.Duration(-c.Sync_interval) * time.Second)
		} else {
			log.Debug("SYNC COUNT eq 0")
			gte.Time = lte.Add(time.Duration(-c.Sync_depth) * time.Minute)

		}

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
						"query":            string(service),
						"fields":           []string{"message"},
						"analyze_wildcard": bool(true),
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
			Size:  500,
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
		uri_ary := []string{"http://", c.Elasticsearch_url, ":", c.Elasticsearch_port, "/", "/_search?pretty"} //c.Elasticsearch_index, "/_search?pretty"}
		query_uri := strings.Join(uri_ary, "")
		log.Debug("Query URI: ", query_uri)
		// Make request
		req, err := http.NewRequest("POST", query_uri, bytes.NewBuffer(jsonpost))
		if err != nil {
			log.Error(err)
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
		}

		// Print
		for k0, v0 := range response.Hits.(map[string]interface{}) {
			if k0 == "hits" {
				for _, v1 := range v0.([]interface{}) {
					for k2, v2 := range v1.(map[string]interface{}) {
						if k2 == "_source" {
							log.Debug("Source: ", v2)
							for key, message := range v2.(map[string]interface{}) {
								if key == "message" {
									log.Info(message.(string))
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
		query := *define_service
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
	fmt.Println(`██╗      ██████╗  ██████╗ ██╗████████╗`)
	fmt.Println(`██║     ██╔═══██╗██╔════╝ ██║╚══██╔══╝`)
	fmt.Println(`██║     ██║   ██║██║  ███╗██║   ██║   `)
	fmt.Println(`██║     ██║   ██║██║   ██║██║   ██║   `)
	fmt.Println(`███████╗╚██████╔╝╚██████╔╝██║   ██║   `)
	fmt.Println(`╚══════╝ ╚═════╝  ╚═════╝ ╚═╝   ╚═╝   `)
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

	config, err := options(*config_path)
	if err != nil {
		log.Error(err)
	}

	log.Debug("Configuration:")
	log.Debug(
		"Defines: ", config.Define, "\n",
		"Sync: ", config.Sync_interval, "\n",
		"ES URL: ", config.Elasticsearch_url, "\n",
		"ES Port: ", config.Elasticsearch_port, "\n",
		"ES Index: ", config.Elasticsearch_index)

	defines := config.Define
	svc_query := lookup(defines)
	log.Info("Querying ", *service, ": ", svc_query)
	// Roll into the query loop
	query(svc_query, config)
}
