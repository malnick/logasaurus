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
)

// Config struct
type Config struct {
	Define              map[string]string `yaml:"define"`
	Sync_interval       int               `yaml:"sync_interval"`
	Elasticsearch_url   string            `yaml:"elasticsearch_url"`
	Elasticsearch_port  string            `yaml:"elasticsearch_port"`
	Elasticsearch_index string            `yaml:"elasticsearch_index"`
}

type Es_resp struct {
	Hits interface{}
	//Hits map["hits"]map["hits"][]map["_source"]map["message"]string
}

type Es_post struct {
	Size  int                          `json:"size"`
	Sort  map[string]map[string]string `json:"sort"`
	Query map[string]interface{}       `json:"query"` //map[string]map[string]map[string]map[string]interface{} `json:"query"`
}

// Define flag overrides
var config_path = flag.String("c", "./config.yaml", "The path to the logit.yaml. Default: ~/.logit/config.yaml (osx) and /etc/logit/config.yaml (*nix).")
var define_service = flag.String("d", "", "A one-time defined service. Must be valid ES query.")
var elastic_url = flag.String("e", "", "Elastic search URL. Default: localhost:9300")
var sync_interval = flag.Int("i", 0, "Query interval in seconds. Default: 1")
var elastic_port = flag.String("p", "", "Elastic Search port. Default: 9200.")
var elastic_index = flag.String("in", "", "Elastic Search index. Default: logstash-\\*")
var verbose = flag.Bool("v", false, "Verbosity. Default: false")
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

	return o, nil
}

func query(service string, c Config) (response Es_resp, err error) {
	// The JSON
	sort := map[string]map[string]string{
		"@timestamp": map[string]string{
			"order":         "desc",
			"unmapped_type": "true",
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
			"filter": map[string]map[string][]map[string]map[string]map[string]string{
				"bool": {
					"must": {
						{
							"range": {
								"@timestamp": {
									"gte": string("1440698405782"),
									"lte": string("1440699305782"),
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
		Size:  5,
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
		return response, err
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
		return response, err
	}
	return response, nil
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

	// Query string: from CLI or config.yaml?
	// Assert map string string type for defines sub map in config
	defines := config.Define
	svc_query := lookup(defines)
	log.Info("Querying ", *service, ": ", svc_query)

	full_response, _ := query(svc_query, config)
	log.Info(full_response)
}
