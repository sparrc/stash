package stash

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

// Test loading the test config file
func TestLoad(t *testing.T) {
	expectConfig := ConfigEntries{
		testConfigEntry(),
	}
	config := getTestConfig()
	if !reflect.DeepEqual(config.Entries, expectConfig) {
		t.Errorf("EXPECTED %s GOT %s",
			expectConfig,
			config.Entries)
	}
}

// Test that function properly loads previous configs and adds new config
func TestAdd(t *testing.T) {
	var t0 time.Time
	newEntry := ConfigEntry{
		Name:        "Wahoo",
		Folders:     []string{"/home"},
		Type:        "Google",
		Frequency:   time.Duration(0),
		LastBak:     t0,
		Credentials: map[string]string{"apikey": "12345"},
	}
	expectConfig := ConfigEntries{
		testConfigEntry(),
		newEntry,
	}
	configMngr := getTestConfig()
	newConfig := configMngr.getNewConfigEntries(newEntry)
	if !reflect.DeepEqual(newConfig, expectConfig) {
		t.Errorf("EXPECTED %s GOT %s",
			expectConfig,
			newConfig)
	}
}

// Test JSON marshalling a config entry
func TestJSONMarshall(t *testing.T) {
	testConfig := ConfigEntries{
		testConfigEntry(),
	}
	expectStr := `[
  {
    "Name": "FooBar",
    "Folders": [
      "/tmp/foo",
      "/tmp/bar"
    ],
    "Type": "Amazon",
    "Frequency": 0,
    "LastBak": "0001-01-01T00:00:00Z",
    "Credentials": {
      "key": "supersecret",
      "keyID": "123"
    }
  }
]
`
	configMngr := getTestConfig()
	testStr := configMngr.ToJSON(testConfig)
	if strings.Trim(string(testStr), " \n") != strings.Trim(expectStr, " \n") {
		t.Errorf("\nEXPECTED\n%s\nGOT\n%s",
			expectStr,
			testStr)
	}
}

// Test that duplicate configs get filtered out
func TestAddDuplicate(t *testing.T) {
	configMngr := getTestConfig()
	newEntries := configMngr.getNewConfigEntries(testConfigEntry())
	if len(newEntries) > 1 {
		t.Errorf("Duplicate entry was not properly filtered, config file:\n%s",
			newEntries)
	}
}

func testConfigEntry() ConfigEntry {
	t := time.Date(0001, time.January, 01, 0, 0, 0, 0, time.UTC)
	return ConfigEntry{
		Name:        "FooBar",
		Folders:     []string{"/tmp/foo", "/tmp/bar"},
		Type:        "Amazon",
		Frequency:   time.Duration(0),
		LastBak:     t,
		Credentials: map[string]string{"key": "supersecret", "keyID": "123"},
	}
}

func getTestConfig() *Config {
	wd, _ := os.Getwd()
	filename := filepath.Join(wd, "config_test.json")
	config := Config{
		FileName: filename,
		Entries:  loadConfig(filename),
	}
	return &config
}
