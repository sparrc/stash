package config

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func getConfigEntry() Entry {
	return Entry{
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
	expectConfig := []Entry{
		getConfigEntry(),
	}
	configMngr := Mngr{FileName: testConfFile}
	fileConfig := configMngr.LoadConfigFile()
	if !reflect.DeepEqual(fileConfig, expectConfig) {
		t.Errorf("EXPECTED %s GOT %s",
			expectConfig,
			fileConfig)
	}
}

// Test that function properly loads previous configs and adds new config
func TestAdd(t *testing.T) {
	testConfFile := getConfFile()
	newEntry := Entry{
		Name:        "Wahoo",
		Folders:     []string{"/home"},
		Type:        "Google",
		Credentials: map[string]string{"apikey": "12345"},
	}
	expectConfig := []Entry{
		getConfigEntry(),
		newEntry,
	}
	configMngr := Mngr{FileName: testConfFile}
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
	testConfig := []Entry{
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
    "Credentials": {
      "key": "supersecret",
      "keyID": "123"
    }
  }
]
`
	configMngr := Mngr{FileName: testConfFile}
	testStr := configMngr.JSONMarshallEntry(testConfig)
	if strings.Trim(string(testStr), " \n") != strings.Trim(expectStr, " \n") {
		t.Errorf("\nEXPECTED\n%s\nGOT\n%s",
			expectStr,
			testStr)
	}
}

// Test adding duplicate configs
