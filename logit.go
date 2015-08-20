package main

import (
	"log"
	"flag"
)

type Flags struct {
	config bool,
	define string,
	elastic_uri string,
	interval int,
	index string,
	port int,
	verbose bool,
}

func options() (f Flags, err error)  {
	f.Config := flag.Bool("config", false, "The path to the logit.yaml. Default: ~/.logit/config.yaml (osx) and /etc/logit/config.yaml (*nix).")

	return f
}

func main() {
	// Get cli flags
	flags := options()

}
