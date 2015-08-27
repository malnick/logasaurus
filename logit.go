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

func options() map[string]interface{} {
	// Create a map to drop arbitrary values into
	o := make(map[string]interface{})
	// Add values from flags to map
	o["config"] = flag.String("c", "~/.logit/config.yaml", "The path to the logit.yaml. Default: ~/.logit/config.yaml (osx) and /etc/logit/config.yaml (*nix).")
	o["define"] = flag.String("d", "", "A one-time defined service. Must be valid ES query.")
	o["elastic_url"] = flag.String("e", "localhost:9300", "Elastic search URL. Default: localhost:9300")
	o["interval"] = flag.Int("i", 1, "Query interval in seconds. Default: 1")
	return o
}

func main() {
	// Get cli flags
	flags := options()
	flag.Parse()
	log.Println(flags["define"])
	log.Println("Configuration:")
	for k, v := range flags {
		log.Println(k, v)
	}
}
