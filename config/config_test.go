package config

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	testConfig := defaultConfig()
	expected := Config{
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

	if testConfig.SyncInterval != expected.SyncInterval {
		t.Error("Expected default sync interval to be 5, got", testConfig.SyncInterval)
	}

	if testConfig.SyncDepth != expected.SyncDepth {
		t.Error("Expected default sync depth to be 10, got", testConfig.SyncDepth)
	}

	if testConfig.ElasticsearchPort != expected.ElasticsearchPort {
		t.Error("Expected default ES port to be 9200, got", testConfig.ElasticsearchPort)
	}

	if testConfig.SearchHost != expected.SearchHost {
		t.Error("Expected default search host to be true, got", testConfig.SearchHost)
	}

	if testConfig.Highlight != expected.Highlight {
		t.Error("Expected default highlight to be true, got", testConfig.Highlight)
	}

	if testConfig.StartTime != expected.StartTime {
		t.Error("Expected default start time to be 0, got", testConfig.StartTime)
	}

	if testConfig.Count != expected.Count {
		t.Error("Expected default count to be 500, got", testConfig.Count)
	}

	if testConfig.LogVerbose != expected.LogVerbose {
		t.Error("Expected default log verbose to be false, got", testConfig.LogVerbose)
	}

	if testConfig.logaConfigPath != expected.logaConfigPath {
		t.Error("Expected default config path to be", expected.logaConfigPath, "got", testConfig.logaConfigPath)
	}
}

func TestFromLogaYaml(t *testing.T) {
	var config = defaultConfig()
	var (
		yamlConfig = `
sync_interval: 10
count: 20
`
		badConfig = `
 bar
`
	)
	file, err := ioutil.TempFile(os.TempDir(), "loga_yaml")
	if err != nil {
		t.Error("Could not make temp file")
	}
	defer os.Remove(file.Name())
	config.logaConfigPath = file.Name()
	file.Write([]byte(yamlConfig))
	err = config.fromLogaYaml()
	if err != nil {
		t.Error("Count not get configuration from temp file")
	}

	if config.SyncInterval != 10 {
		t.Error("Expected sync interval from config to be 10, got", config.SyncInterval)
	}

	if config.Count != 20 {
		t.Error("Expected count from config to be 20, got", config.Count)
	}

	file.Write([]byte(badConfig))
	err = config.fromLogaYaml()
	if err == nil {
		t.Error("expected error, got", err)
	}
}

func TestGetDefinedQuery(t *testing.T) {
	var (
		config     = Config{}
		yamlConfig = `
define_service:
  test: "foo AND bar"
`
	)
	file, err := ioutil.TempFile(os.TempDir(), "loga_yaml")
	if err != nil {
		t.Error(err)
	}

	defer os.Remove(file.Name())
	config.logaConfigPath = file.Name()
	file.Write([]byte(yamlConfig))
	if err := config.fromLogaYaml(); err != nil {
		t.Error(err)
	}

	config.FlagDefinedQuery = "bar AND foo"
	query, err := config.GetDefinedQuery()
	if err != nil {
		t.Error("expected no errors getting query, got", err)
	}
	if query != "bar AND foo" {
		t.Error("expected query to be 'bar AND foo', got", query)
	}

	config.FlagDefinedQuery = ""
	config.FlagConfDefinedQuery = "test"
	fooQuery, err := config.GetDefinedQuery()
	if err != nil {
		t.Error("expected no errors getting query, got", err)
	}
	if fooQuery != "foo AND bar" {
		t.Error("expected query to be 'foo AND bar', got", fooQuery)
	}
}
