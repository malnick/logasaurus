package config

import "testing"

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
