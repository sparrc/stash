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

// ConfigEntries is an array of configuration entries
type ConfigEntries []ConfigEntry

// NewConfig creates a new configuration manager with default file path set.
func NewConfig() *Config {

	stashdir := filepath.Join(os.Getenv("HOME"), ".stash")
	// Get the config file location path
	filename := filepath.Join(stashdir, "config")

	// If DB directory does not exist, create it
	if _, err := os.Stat(stashdir); os.IsNotExist(err) {
		exec.Command("mkdir", "-p", stashdir).Run()
	}

	// Create the config struct
	config := Config{
		FileName: filename,
		Entries:  loadConfig(filename),
	}

	return &config
}

// AddDestination adds a backup destination to the config file
func (cm *Config) AddDestination(configEntry ConfigEntry) error {
	log.Printf("Adding destination [%s] to config file [%s]\n",
		configEntry.Name,
		cm.FileName)

	// TODO handle timeout if someone else has it open
	db, err := bolt.Open(cm.FileName, 0666, nil)
	if err != nil {
		log.Fatal("Error opening the database file", err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("destinations"))
		if err != nil {
			return err
		}

		data, err := json.Marshal(configEntry)
		if err != nil {
			return err
		}

		b.Put([]byte(configEntry.Name), data)

		return nil
	})
	return err
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

// TouchLastBak updates the "LastBak" timestamp to time.Now()
func (cm *Config) TouchLastBak(name string) error {
	// TODO handle timeout if someone else has it open
	db, err := bolt.Open(cm.FileName, 0666, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("destinations"))

		// bucket doesnt exist:
		if b == nil {
			return nil
		}

		// Get the entry and set new LastBak time
		v := b.Get([]byte(name))
		var ce ConfigEntry
		json.Unmarshal(v, &ce)
		// Subtract a millisecond to avoid polling creep
		ce.LastBak = time.Now().Add(-1 * time.Millisecond)

		// Put the updated entry back into the DB
		data, err := json.Marshal(ce)
		if err != nil {
			return err
		}
		b.Put([]byte(name), data)

		return nil
	})
	return err
}

// DeleteEntry deletes the given entry
func (cm *Config) DeleteEntry(name string) error {
	// TODO handle timeout if someone else has it open
	db, err := bolt.Open(cm.FileName, 0666, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("destinations"))

		// bucket doesnt exist, so do nothing:
		if b == nil {
			return nil
		}

		// Delete the entry (this returns nil if entry doesnt exist)
		err := b.Delete([]byte(name))

		return err
	})
	return err
}

// loadConfig loads the given config db file and returns ConfigEntries.
// In general, if anything goes wrong, it just returns empty ConfigEntries.
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
		log.Println("Error loading config: ", err)
	}
	return entries
}
