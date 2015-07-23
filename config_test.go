package stash

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

// Test loading the test config file
func TestLoad(t *testing.T) {
	t0 := time.Date(0001, time.January, 01, 0, 0, 0, 0, time.UTC)
	// This config matches that in the testdata/config file
	expectConfig := ConfigEntries{
		ConfigEntry{
			Name:        "FooBar",
			Folders:     []string{"/tmp"},
			Type:        "Amazon",
			Frequency:   time.Duration(0),
			LastBak:     t0,
			Credentials: map[string]string{"key": "supersecret"},
		},
	}
	config := loadTestConfig()
	if !reflect.DeepEqual(config.Entries, expectConfig) {
		t.Errorf("EXPECTED %s GOT %s",
			expectConfig,
			config.Entries)
	}
}

// Test IsDuplicateEntry returns true/false correctly
func TestIsDuplicateEntry(t *testing.T) {
	t0 := time.Date(0001, time.January, 01, 0, 0, 0, 0, time.UTC)
	dupeConfig := ConfigEntry{
		Name:        "FooBar",
		Folders:     []string{"/tmp"},
		Type:        "Amazon",
		Frequency:   time.Duration(0),
		LastBak:     t0,
		Credentials: map[string]string{"key": "supersecret"},
	}
	config := loadTestConfig()
	if !config.IsDuplicateEntry(dupeConfig) {
		t.Error("Expected IsDuplicateEntry to return True, dupeEntry: ",
			dupeConfig, " fileEntry: ", config.Entries)
	}

	// Change dupeConfig and verify that IsDuplicateEntry returns False:
	dupeConfig.Name = "NotDupe"
	if config.IsDuplicateEntry(dupeConfig) {
		t.Error("Expected IsDuplicateEntry to return False, dupeEntry: ",
			dupeConfig, " fileEntry: ", config.Entries)
	}
}

// Test that we can add a new config entry, load, and read it back identically
func TestAddReload(t *testing.T) {
	tmp := tempdb()
	defer os.Remove(tmp)
	config := Config{
		FileName: tmp,
		Entries:  loadConfig(tmp),
	}

	t0 := time.Date(0001, time.January, 01, 0, 0, 0, 0, time.UTC)
	newEntry := ConfigEntry{
		Name:        "Wahoo",
		Folders:     []string{"/home"},
		Type:        "Google",
		Frequency:   time.Duration(0),
		LastBak:     t0,
		Credentials: map[string]string{"apikey": "12345"},
	}

	config.AddDestination(newEntry)
	config.ReloadConfig()
	if !reflect.DeepEqual(config.Entries[0], newEntry) {
		t.Errorf("EXPECTED %s GOT %s",
			newEntry,
			config.Entries[0])
	}
}

func loadTestConfig() *Config {
	wd, _ := os.Getwd()
	filename := filepath.Join(wd, "testdata", "config_test")
	config := Config{
		FileName: filename,
		Entries:  loadConfig(filename),
	}
	return &config
}

func tempdb() string {
	wd, _ := os.Getwd()
	dir := filepath.Join(wd, "testdata")
	f, _ := ioutil.TempFile(dir, "tmpdb")
	fname := f.Name()
	defer f.Close()
	defer os.Remove(fname)
	return fname
}
