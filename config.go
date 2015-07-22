package stash

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
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
	filename := filepath.Join(os.Getenv("HOME"), ".stash", "config.json")

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
	JSON := cm.ToJSON(allConfigEntries)
	ioutil.WriteFile(cm.FileName, JSON, 0644)
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

// ToJSON marshalls ConfigEntries into JSON
func (cm *Config) ToJSON(configEntries ConfigEntries) []byte {
	JSON, err := json.MarshalIndent(configEntries, "", "  ")
	if err != nil {
		panic(err)
	}
	return JSON
}

// LoadConfig loads the config file and returns the contents
func (cm *Config) ReloadConfig() {
	cm.Entries = loadConfig(cm.FileName)
}

func loadConfig(filename string) ConfigEntries {
	log.Printf("Loading config file [%s]\n", filename)
	content, err := ioutil.ReadFile(filename)
	var entries ConfigEntries
	if err != nil {
		// This indicates there was no config file present, return empty entries
		return entries
	} else if len(content) == 0 {
		return entries
	}

	if err := json.Unmarshal(content, &entries); err != nil {
		log.Println("Error loading config: ", err)
	}
	return entries
}
