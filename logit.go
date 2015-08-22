package main

import (
	"flag"
	"log"
)

//	config      *bool
//	define      *string
//	elastic_uri *string
//	interval    *int
//	index       *string
//	port        *int
//	verbose     *bool
//}

var config_file map[string]string

func options() (o map[string]interface{}) {

	config := flag.Bool("c", false, "The path to the logit.yaml. Default: ~/.logit/config.yaml (osx) and /etc/logit/config.yaml (*nix).")
	o["config"] = config

	//	o["define"] = flag.String("d", "", "A one-time defined service. Must be valid ES query.")
	//	o["elastic_url"] = flag.String("e", "localhost:9300", "Elastic search URL. Default: localhost:9300")
	//	o["interval"] = flag.Int("i", 1, "Query interval in seconds. Default: 1")

	return o
}

func main() {
	// Get cli flags
	flags := options()
	log.Println("Configuration:\n")
	for k, v := range flags {
		log.Println(k, v)
	}
}
