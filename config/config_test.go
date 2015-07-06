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
			expectConfig,
			fileConfig)
	}
}

// Test adding another destination and loading it back
func TestAdd(t *testing.T) {
	wd, _ := os.Getwd()
	testConfFile := filepath.Join(wd, "config_test.json")
	newEntry := ConfigEntry{
		Name:        "Wahoo",
		Folders:     []string{"/home"},
		Type:        "Google",
		Credentials: map[string]string{"apikey": "12345"},
	}
	expectConfig := []ConfigEntry{
		ConfigEntry{
			Name:        "FooBar",
			Folders:     []string{"/tmp/foo", "/tmp/bar"},
			Type:        "Amazon",
			Credentials: map[string]string{"key": "supersecret", "keyID": "123"},
		},
		newEntry,
	}
	configMngr := ConfigMngr{FileName: testConfFile}
	newConfig := configMngr.GetNewConfigEntries(newEntry)
	if !reflect.DeepEqual(newConfig, expectConfig) {
		t.Errorf("EXPECTED %s GOT %s",
			expectConfig,
			newConfig)
	}
}

// Test adding duplicate configs
