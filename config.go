package main

import (
	"flag"
	"time"
)

type Config struct {
	Define             map[string]string `yaml:"define"`
	SyncInterval       int               `yaml:"sync_interval"`
	SyncDepth          int               `yaml:"sync_depth"`
	ElasticsearchURL   string            `yaml:"elasticsearch_url"`
	ElasticsearchPort  string            `yaml:"elasticsearch_port"`
	ElasticsearchIndex string            `yaml:"elasticsearch_index"`
	SearchHost         bool
	Highlight          bool
	StartTime          time.Time
	Count              int
	LogVerbose         bool
	logaConfigPath     string
}

func defaultConfig() Config {
	return Config{
		SyncInterval:       5,
		SyncDepth:          10,
		ElasticsearchURL:   "",
		ElasticsearchPort:  "9200",
		ElasticsearchIndex: "",
		SearchHost:         false,
		Highlight:          true,
		StartTime:          time.Now(),
		Count:              500,
		LogVerbose:         false,
		logaConfigPath:     "./loga.yaml",
	}

}

func (c *Config) setFlags(fs *flag.FlagSet) {
	fs.BoolVar(&c.LogVerbose, "v", c.LogVerbose, "Verbose logging option")
}

func ParseArgsReturnConfig() Config {
	config := defaultConfig()
	logaFlags := flag.NewFlagSet("", flag.ContinueOnError)
	config.setFlags(logaFlags)
	return config
}
