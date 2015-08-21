package main

import (
	"flag"
	"log"
)

type Flags struct {
	config      *bool
	define      *string
	elastic_uri *string
	interval    *int
	index       *string
	port        *int
	verbose     *bool
}

func options() (f Flags) {
	config := flag.Bool("config", false, "The path to the logit.yaml. Default: ~/.logit/config.yaml (osx) and /etc/logit/config.yaml (*nix).")
	define := flag.String("define", "", "A one-time defined service. Must be valid ES query.")

	f.config = config
	f.define = define
	return f
}

func main() {
	// Get cli flags
	flags := options()
	log.Println("Configuration:\n")
	for k, v := range flags {
		log.Println(k, v)
	}
}
