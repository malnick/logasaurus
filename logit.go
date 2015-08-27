package main

import (
	"flag"
	"log"
)

// Define flag overrides
var config_path = flag.String("c", "~/.logit/config.yaml", "The path to the logit.yaml. Default: ~/.logit/config.yaml (osx) and /etc/logit/config.yaml (*nix).")
var define_service = flag.String("d", "", "A one-time defined service. Must be valid ES query.")
var elastic_url = flag.String("e", "localhost:9300", "Elastic search URL. Default: localhost:9300")
var sync_interval = flag.Int("i", 1, "Query interval in seconds. Default: 1")

func options() (o map[string]interface{}, err error) {

	return o, nil
}

func main() {
	// Get cli flags
	flag.Parse()
	log.Println("Define: ", *define_service)
}
