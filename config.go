package stash

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
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

	// Create config file if it doesn't exist:
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		config.createConfig()
	}
	return &config
}

func (cm *Config) createConfig() {
	log.Println("Creating config file")
	args := []string{"-p", filepath.Dir(cm.FileName)}
	if err := exec.Command("mkdir", args...).Run(); err != nil {
		panic(err)
	}
	if err := exec.Command("touch", cm.FileName).Run(); err != nil {
		panic(err)
	}
}

// AddDestination adds a backup destination to the config file
func (cm *Config) AddDestination(configEntry ConfigEntry) {
	log.Printf("Adding destination [%s] to config file [%s]\n",
		configEntry.Name,
		cm.FileName)

	allConfigEntries := cm.getNewConfigEntries(configEntry)
	// TODO handle timeout if someone else has it open
	db, err := bolt.Open(cm.FileName, 0666, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err := db.Update(func(tx *bolt.Tx) error {
		// TODO handle error
		b, _ := tx.CreateBucketIfNotExists([]byte("destinations"))

		for _, entry := range allConfigEntries {
			k := entry.Name
			// TODO handle error
			data, _ := json.Marshal(entry)
			b.Put([]byte(k), data)
		}

		return nil
	}); err != nil {
		panic(err)
	}
}

// getNewConfigEntries takes a new config entry, loads previous entries,
// combines & removes duplicate entries.
func (cm *Config) getNewConfigEntries(newEntry ConfigEntry) ConfigEntries {
	cm.ReloadConfig()
	if !cm.IsDuplicateEntry(newEntry) {
		cm.Entries = append(cm.Entries, newEntry)
	}
	return cm.Entries
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

// LoadConfig loads the config file and returns the contents
func (cm *Config) ReloadConfig() {
	cm.Entries = loadConfig(cm.FileName)
}

func loadConfig(filename string) ConfigEntries {
	log.Printf("Loading config file [%s]\n", filename)

	var entries ConfigEntries

	// TODO: Handle a timeout if someone has it open "rw"
	db, err := bolt.Open(filename, 0666, &bolt.Options{ReadOnly: true})
	if err != nil {
		// Indicates that there are currently no entries if there is no db
		return entries
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
