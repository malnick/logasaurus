package main

import (
	"flag"
	"os"
	"strings"
	"time"
)

type Config struct {
	Define             map[string]string `yaml:"define"`
	SyncInterval       int               `yaml:"sync_interval"`
	SyncDepth          int               `yaml:"sync_depth"`
	ElasticsearchURL   string            `yaml:"elasticsearch_url"`
	ElasticsearchPort  string            `yaml:"elasticsearch_port"`
	ElasticsearchIndex string            `yaml:"elasticsearch_index"`
	host               bool
	highlight          bool
	startTime          time.Time
	count              int
	logaHome           string
	logaConfigPath     string
}

func defaultConfig() Config {
	return Config{
		SyncInterval: 5,
		SyncDepth:    10,

		logaHome: os.Getenv("HOME"),
	}

}

func Configuration() {
	var c Config
	c.logaConfigPath = flag.String("c", strings.Join([]string{home, "/.loga/config.yaml"}, ""), "The path to the config.yaml.")

}
