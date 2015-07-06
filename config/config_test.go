package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// Test loading the test config file
func TestLoad(t *testing.T) {
	wd, _ := os.Getwd()
	testConfFile := filepath.Join(wd, "config_test.json")
	expectConfig := []ConfigEntry{
		ConfigEntry{
			Name:        "FooBar",
			Folders:     []string{"/tmp/foo", "/tmp/bar"},
			Type:        "Amazon",
			Credentials: map[string]string{"key": "supersecret", "keyID": "123"},
		},
	}
	configMngr := ConfigMngr{FileName: testConfFile}
	fileConfig := configMngr.LoadConfigFile()
	if !reflect.DeepEqual(fileConfig, expectConfig) {
		t.Errorf("EXPECTED %s GOT %s",
			expectConfig[0],
			fileConfig[0])
	}
}

// Test adding another destination and loading it back

// Test adding duplicate configs
