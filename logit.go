package main

import (
	"log"
	"flag"
)

type Flags struct {
	config bool
	define string
	elastic_uri string
	interval int
	index string
	port int
	verbose bool
}

func options() Flags {
	f.Config := flag.Bool("config", false, "The path to the logit.yaml. Default: ~/.logit/config.yaml (osx) and /etc/logit/config.yaml (*nix).")
	log.Println("Configuration:")
	log.Println(f)
	return f, nil
}

func main() {
	// Get cli flags
	flags := options()
	log.Println(flags)
}
