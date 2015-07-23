package stash

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/boltdb/bolt"
)

// Config manages the stash configuration file.
type Config struct {
	// Name of the configuration file.
	FileName string

	// Entries in the config file
	Entries ConfigEntries
}

// ConfigEntry specifies a configuration entry
type ConfigEntry struct {
	Name        string
	Folders     []string
	Type        string
	Frequency   time.Duration
	LastBak     time.Time
	Credentials map[string]string
}

type ConfigEntries []ConfigEntry

// NewConfig creates a new configuration manager with default file path set.
func NewConfig() *Config {
	// Get the config file location path
	filename := filepath.Join(os.Getenv("HOME"), ".stash", "config")

	// Create the config struct
	config := Config{
		FileName: filename,
		Entries:  loadConfig(filename),
	}

	return &config
}

// AddDestination adds a backup destination to the config file
func (cm *Config) AddDestination(configEntry ConfigEntry) {
	log.Printf("Adding destination [%s] to config file [%s]\n",
		configEntry.Name,
		cm.FileName)

	// TODO handle timeout if someone else has it open
	db, err := bolt.Open(cm.FileName, 0666, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err := db.Update(func(tx *bolt.Tx) error {
		// TODO handle error
		b, _ := tx.CreateBucketIfNotExists([]byte("destinations"))

		// TODO handle error
		data, _ := json.Marshal(configEntry)
		b.Put([]byte(configEntry.Name), data)

		return nil
	}); err != nil {
		panic(err)
	}
}

// IsDuplicateEntry returns true if the entry already exists in the config file
func (cm *Config) IsDuplicateEntry(newEntry ConfigEntry) bool {
	for _, entry := range cm.Entries {
		if entry.Name == newEntry.Name {
			return true
		}
	}
	return false
}

// ReloadConfig loads the config file and returns the contents
func (cm *Config) ReloadConfig() {
	cm.Entries = loadConfig(cm.FileName)
}

func loadConfig(filename string) ConfigEntries {
	log.Printf("Loading config file [%s]\n", filename)

	var entries ConfigEntries

	// If DB file does not exist, just return empty entries
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return entries
	}

	// TODO: Handle a db open timeout if someone has it open "rw"
	db, err := bolt.Open(filename, 0666, &bolt.Options{ReadOnly: true})
	if err != nil {
		// Something went terribly wrong
		panic(err)
	}
	defer db.Close()

	if err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("destinations"))

		// Bucket doesnt exist:
		if b == nil {
			return nil
		}

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var ce ConfigEntry
			json.Unmarshal(v, &ce)
			entries = append(entries, ce)
		}

		return nil
	}); err != nil {
		panic(err)
	}
	return entries
}
