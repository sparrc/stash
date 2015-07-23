package main

import (
	"log"
	"time"

	"github.com/sparrc/stash"
)

func main() {
	ticker := time.NewTicker(10 * time.Second)
	quit := make(chan bool)
	go daemon(ticker, quit)
	time.Sleep(120 * time.Second)
	close(quit)
}

func daemon(ticker *time.Ticker, quit <-chan bool) {
	config := stash.NewConfig()
	for {
		select {
		case <-ticker.C:
			processBackups(config)
		case <-quit:
			ticker.Stop()
			return
		}
	}
}

func processBackups(config *stash.Config) {
	for _, entry := range config.Entries {
		// log.Println("Processing Backup:	", entry.Name)
		if entry.LastBak.Add(entry.Frequency).Before(time.Now()) {
			doBackup(entry)
			config.TouchLastBak(entry.Name)
			config.ReloadConfig()
		}
	}
}

func doBackup(entry stash.ConfigEntry) {
	log.Println("Performing Backup:	", entry.Name)
	switch entry.Type {
	case "Amazon":
		doAmazon(entry)
	case "Google":
		doGoogle(entry)
	}
}

func doAmazon(entry stash.ConfigEntry) {
	return
}

func doGoogle(entry stash.ConfigEntry) {
	return
}
