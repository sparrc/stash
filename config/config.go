package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// Mngr manages the stash configuration file.
type Mngr struct {
	// Name of the configuration file.
	FileName string
}

// Entry specifies a configuration entry
type Entry struct {
	Name        string
	Folders     []string
	Type        string
	Credentials map[string]string
}

// NewConfigMngr creates a new configuration manager with default file path set.
func NewConfigMngr() *Mngr {
	filename := filepath.Join(os.Getenv("HOME"), ".stash", "config.json")
	config := Mngr{FileName: filename}
	// Create config file if it doesn't exist:
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		config.createConfigFile()
	}
	return &config
}

func (cm *Mngr) createConfigFile() {
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
func (cm *Mngr) AddDestination(configEntry Entry) {
	fmt.Fprintf(os.Stdout,
		"Adding destination [%s] to config file [%s]\n",
		configEntry.Name,
		cm.FileName)

	allConfigEntries := cm.GetNewConfigEntries(configEntry)
	JSON := cm.JSONMarshallEntry(allConfigEntries)
	ioutil.WriteFile(cm.FileName, JSON, 0644)
}

// GetNewConfigEntries takes a new config entry, loads previous entries,
//					combines & removes duplicate entries.
func (cm *Mngr) GetNewConfigEntries(newEntry Entry) []Entry {
	configFile := cm.LoadConfigFile()
	duplicate := false
	for _, entry := range configFile {
		if entry.Name == newEntry.Name {
			duplicate = true
			fmt.Fprintf(os.Stderr,
				"stash: Attempting to add duplicate entry %s\n",
				newEntry.Name)
		}
	}
	if !duplicate {
		configFile = append(configFile, newEntry)
	}
	return configFile
}

// JSONMarshallEntry marshalls a config.Entry into JSON
func (cm *Mngr) JSONMarshallEntry(configEntries []Entry) []byte {
	JSON, err := json.MarshalIndent(configEntries, "", "  ")
	if err != nil {
		panic(err)
	}
	return JSON
}

// LoadConfigFile loads the config file and returns the contents
func (cm *Mngr) LoadConfigFile() []Entry {
	fmt.Fprintf(os.Stdout, "Loading config file [%s]\n", cm.FileName)
	content, err := ioutil.ReadFile(cm.FileName)
	if err != nil {
		panic(err)
	}
	var entries []Entry
	if err := json.Unmarshal(content, &entries); err != nil {
		log.Println("No config entries loaded.")
	}
	return entries
}
