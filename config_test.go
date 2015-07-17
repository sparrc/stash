package stash

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func getConfigEntry() ConfigEntry {
	return ConfigEntry{
		Name:        "FooBar",
		Folders:     []string{"/tmp/foo", "/tmp/bar"},
		Type:        "Amazon",
		Credentials: map[string]string{"key": "supersecret", "keyID": "123"},
	}
}

func getConfFile() string {
	wd, _ := os.Getwd()
	return filepath.Join(wd, "config_test.json")
}

// Test loading the test config file
func TestLoad(t *testing.T) {
	testConfFile := getConfFile()
	expectConfig := []ConfigEntry{
		getConfigEntry(),
	}
	configMngr := Config{FileName: testConfFile}
	fileConfig := configMngr.LoadConfig()
	if !reflect.DeepEqual(fileConfig, expectConfig) {
		t.Errorf("EXPECTED %s GOT %s",
			expectConfig,
			fileConfig)
	}
}

// Test that function properly loads previous configs and adds new config
func TestAdd(t *testing.T) {
	testConfFile := getConfFile()
	newEntry := ConfigEntry{
		Name:        "Wahoo",
		Folders:     []string{"/home"},
		Type:        "Google",
		Credentials: map[string]string{"apikey": "12345"},
	}
	expectConfig := []ConfigEntry{
		getConfigEntry(),
		newEntry,
	}
	configMngr := Config{FileName: testConfFile}
	newConfig := configMngr.GetNewConfigEntries(newEntry)
	if !reflect.DeepEqual(newConfig, expectConfig) {
		t.Errorf("EXPECTED %s GOT %s",
			expectConfig,
			newConfig)
	}
}

// Test JSON marshalling a config entry
func TestJSONMarshall(t *testing.T) {
	testConfFile := getConfFile()
	testConfig := []ConfigEntry{
		getConfigEntry(),
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
    "Credentials": {
      "key": "supersecret",
      "keyID": "123"
    }
  }
]
`
	configMngr := Config{FileName: testConfFile}
	testStr := configMngr.ToJSON(testConfig)
	if strings.Trim(string(testStr), " \n") != strings.Trim(expectStr, " \n") {
		t.Errorf("\nEXPECTED\n%s\nGOT\n%s",
			expectStr,
			testStr)
	}
}

// Test that duplicate configs get filtered out
func TestAddDuplicate(t *testing.T) {
	configMngr := Config{FileName: getConfFile()}
	newEntries := configMngr.GetNewConfigEntries(getConfigEntry())
	if len(newEntries) > 1 {
		t.Errorf("Duplicate entry was not properly filtered, config file:\n%s",
			newEntries)
	}
}
