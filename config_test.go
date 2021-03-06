package stash

import (
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
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

// Test that LastBak updates the timestamp to time.Now()
func TestTouchLastBak(t *testing.T) {
	tmp := tempdb()
	defer os.Remove(tmp)
	config := Config{
		FileName: tmp,
		Entries:  loadConfig(tmp),
	}

	t0 := time.Date(0001, time.January, 01, 0, 0, 0, 0, time.UTC)
	newEntry := ConfigEntry{
		Name:    "Wahoo",
		LastBak: t0,
	}

	// Add entry to db
	config.AddDestination(newEntry)

	// Update last backup timestamp
	config.TouchLastBak("Wahoo")
	config.ReloadConfig()

	// Verify that TouchLastBak on a non-existent entry causes no harm:
	config.TouchLastBak("I dont exist")

	// Verify last backup timestamp is past t0
	if !config.Entries[0].LastBak.After(t0) {
		t.Error("Expected last backup date to be past t0\nt0:", t0,
			"Last backup:", config.Entries[0].LastBak)
	}

	// And verify that the timestamp is within 2 seconds:
	if !config.Entries[0].LastBak.After(time.Now().Add(-2 * time.Second)) {
		t.Error("Expected last backup date to be very recent\n",
			"Last backup:", config.Entries[0].LastBak)
	}
}

// Test deleting entries and that deleting non-existent entries does nothing.
func TestDelete(t *testing.T) {
	tmp := tempdb()
	defer os.Remove(tmp)
	config := Config{
		FileName: tmp,
		Entries:  loadConfig(tmp),
	}

	newEntry := ConfigEntry{
		Name: "Wahoo",
	}

	// Add entry to db
	config.AddDestination(newEntry)

	// Delete non-existent entry does nothing:
	err := config.DeleteEntry("I dont exist")
	if err != nil {
		t.Error("Error deleting non-existent entry: ", err)
	}

	// Delete "Wahoo" entry:
	err = config.DeleteEntry("Wahoo")
	if err != nil {
		t.Error("Error deleting Wahoo entry: ", err)
	}

	// Verify that "Wahoo" entry is gone:
	config.ReloadConfig()
	if len(config.Entries) > 0 {
		t.Error("Deleted 'Wahoo' entry but found entries afterwards: ",
			config.Entries)
	}
}

// loads the static test database
func loadTestConfig() *Config {
	wd, _ := os.Getwd()
	filename := filepath.Join(wd, "testdata", "config_test")
	config := Config{
		FileName: filename,
		Entries:  loadConfig(filename),
	}
	return &config
}

// Returns a tmp database filename
func tempdb() string {
	wd, _ := os.Getwd()
	return filepath.Join(wd, "testdata",
		"tmpdb"+strconv.Itoa(rand.Intn(99999)))
}
