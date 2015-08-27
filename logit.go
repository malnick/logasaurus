package main

import (
	"flag"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

// Define flag overrides
var config_path = flag.String("c", "./config.yaml", "The path to the logit.yaml. Default: ~/.logit/config.yaml (osx) and /etc/logit/config.yaml (*nix).")
var define_service = flag.String("d", "", "A one-time defined service. Must be valid ES query.")
var elastic_url = flag.String("e", "", "Elastic search URL. Default: localhost:9300")
var sync_interval = flag.Int("i", 1, "Query interval in seconds. Default: 1")

func options(config_path string) (o map[string]interface{}, err error) {
	config_file, err := ioutil.ReadFile(config_path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(config_file, &o)
	if err != nil {
		return nil, err
	}

	if len(*elastic_url) > 1 {
		o["elasticsearch_url"] = *elastic_url
	}
	return o, nil
}

func main() {
	// Get cli flags
	flag.Parse()
	config, err := options(*config_path)
	if err != nil {
		log.Println(err)
	}
	for k, v := range config {
		log.Println(k, v)
	}
}
