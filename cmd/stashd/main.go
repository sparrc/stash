package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sparrc/stash"
)

func main() {
	// How frequently to poll for backups
	ticker := time.NewTicker(5 * time.Second)

	// Channel to control daemon goroutine
	quit := make(chan bool)
	defer close(quit)
	go daemon(ticker, quit)

	// Run daemon until an interrupt is received
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)
	for {
		select {
		case kill := <-interrupt:
			log.Println("Got signal: ", kill)
			return
		}
	}
}

// daemon controls the backup processors
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
