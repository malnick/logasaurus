package config

import (
	"flag"
	"io/ioutil"

	"gopkg.in/yaml.v2"

	"github.com/malnick/logasaurus/errorhandler"
)

type Config struct {
	DefineService      string
	ConfigService      string `yaml:"define_service"`
	SyncInterval       int    `yaml:"sync_interval"`
	SyncDepth          int    `yaml:"sync_depth"`
	ElasticsearchURL   string `yaml:"elasticsearch_url"`
	ElasticsearchPort  string `yaml:"elasticsearch_port"`
	ElasticsearchIndex string `yaml:"elasticsearch_index"`
	Highlight          bool   `yaml:"highlight"`
	StartTime          int    `yaml:"start_time"`
	Count              int    `yaml:"count"`
	LogVerbose         bool   `yaml:"log_verbose"`
	SearchHost         bool
	logaConfigPath     string
}

func defaultConfig() Config {
	return Config{
		SyncInterval:      5,
		SyncDepth:         10,
		ElasticsearchPort: "9200",
		SearchHost:        false,
		Highlight:         true,
		StartTime:         0,
		Count:             500,
		LogVerbose:        false,
		logaConfigPath:    "./loga.yaml",
	}
}

func (c *Config) fromLogaYaml() {
	configFile, err := ioutil.ReadFile(c.logaConfigPath)
	errorhandler.LogErrorAndExit(err)
	if err := yaml.Unmarshal(configFile, &c); err != nil {
		errorhandler.LogErrorAndExit(err)
	}
}

func (c *Config) setFlags(fs *flag.FlagSet) {
	fs.BoolVar(&c.LogVerbose, "v", c.LogVerbose, "Verbose logging option")
	fs.BoolVar(&c.Highlight, "h", c.Highlight, "Highlight search in output")

	fs.StringVar(&c.logaConfigPath, "c", c.logaConfigPath, "Path to loga.yaml")
	fs.StringVar(&c.DefineService, "d", c.DefineService, "Define a lookup on the CLI")
	fs.StringVar(&c.ConfigService, "s", c.ConfigService, "Name of definition in loga.yaml")
	fs.StringVar(&c.ElasticsearchURL, "e", c.ElasticsearchURL, "URL for Elastic Search")
	fs.StringVar(&c.ElasticsearchPort, "p", c.ElasticsearchPort, "Port for Elastic Search")
	fs.StringVar(&c.ElasticsearchIndex, "in", c.ElasticsearchIndex, "Elastic Search index")

	fs.IntVar(&c.StartTime, "st", c.StartTime, "Start time in minutes. Ex: -st 20 starts query 20 minutes ago.")
	fs.IntVar(&c.SyncInterval, "si", c.SyncInterval, "Query interval in seconds")
	fs.IntVar(&c.SyncDepth, "sd", c.SyncDepth, "Sync depth: how far back to go on initial query")
}

func ParseArgsReturnConfig() Config {
	logaFlags := flag.NewFlagSet("", flag.ContinueOnError)
	config := defaultConfig()
	config.fromLogaYaml()
	config.setFlags(logaFlags)
	return config
}
