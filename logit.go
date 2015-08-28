package main

import (
	//	"encoding/json"
	"bytes"
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

func query(service string) (response Es_resp, err error) {
	var config Config
	// The JSON
	jsonStr := fmt.Sprintf(`{
		"size": 5,
		"sort": [
      {
        "@timestamp": {
          "order": "desc",
          "unmapped_type": "boolean"
        }
      }
    ],
    "query": {
      "filtered": {
        "query": {
          "query_string": {
            "query": "%s",
            "fields": ["message"],
            "analyze_wildcard": true
          }
        },
	}`, service)
	var json = []byte(jsonStr)

	// Craft the request URI
	uri_ary := []string{"http://", config.Elasticsearch_url, ":", config.Elasticsearch_port, "/", config.Elasticsearch_index, "/_search"}
	query_uri := strings.Join(uri_ary, "")
	log.Debug("Query URI: ", query_uri)
	// Make request
	req, err := http.NewRequest("POST", query_uri, bytes.NewBuffer(json))
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

	full_response := query(svc_query)

}
