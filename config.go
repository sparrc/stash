package stash

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// Config manages the stash configuration file.
type Config struct {
	// Name of the configuration file.
	FileName string
}

// ConfigEntry specifies a configuration entry
type ConfigEntry struct {
	Name        string
	Folders     []string
	Type        string
	Credentials map[string]string
}

// NewConfig creates a new configuration manager with default file path set.
func NewConfig() *Config {
	filename := filepath.Join(os.Getenv("HOME"), ".stash", "config.json")
	mngr := Config{FileName: filename}
	// Create config file if it doesn't exist:
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		mngr.createConfig()
	}
	return &mngr
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

	allConfigEntries := cm.GetNewConfigEntries(configEntry)
	JSON := cm.ToJSON(allConfigEntries)
	ioutil.WriteFile(cm.FileName, JSON, 0644)
}

// GetNewConfigEntries takes a new config entry, loads previous entries,
//					combines & removes duplicate entries.
func (cm *Config) GetNewConfigEntries(newEntry ConfigEntry) []ConfigEntry {
	configFile := cm.LoadConfig()
	if !cm.IsDuplicateEntry(newEntry) {
		configFile = append(configFile, newEntry)
	}
	return configFile
}

// IsDuplicateEntry returns true if the entry already exists in the config file
func (cm *Config) IsDuplicateEntry(newEntry ConfigEntry) bool {
	configFile := cm.LoadConfig()
	for _, entry := range configFile {
		if entry.Name == newEntry.Name {
			return true
		}
	}
	return false
}

// ToJSON marshalls a config.ConfigEntry into JSON
func (cm *Config) ToJSON(configEntries []ConfigEntry) []byte {
	JSON, err := json.MarshalIndent(configEntries, "", "  ")
	if err != nil {
		panic(err)
	}
	return JSON
}

// LoadConfig loads the config file and returns the contents
func (cm *Config) LoadConfig() []ConfigEntry {
	log.Printf("Loading config file [%s]\n", cm.FileName)
	content, err := ioutil.ReadFile(cm.FileName)
	if err != nil {
		panic(err)
	}
	var entries []ConfigEntry
	if err := json.Unmarshal(content, &entries); err != nil {
		log.Println("No config entries loaded.")
	}
	return entries
}
