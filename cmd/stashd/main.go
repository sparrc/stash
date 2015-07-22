package main

import (
	"log"
	"time"

	"github.com/sparrc/stash"
)

func main() {
	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan bool)
	go daemon(ticker, quit)
	time.Sleep(30 * time.Second)
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
		log.Println("Processing Backup: ", entry.Name)
		if entry.LastBak.Add(entry.Frequency).Before(time.Now()) {
			log.Println("Performing scheduled backup: ", entry.Name)
		}
	}
}
