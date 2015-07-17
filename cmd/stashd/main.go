package main

import (
	"log"
	"time"

	"github.com/sparrc/stash"
)

func main() {
	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})
	config := stash.NewConfig()
	go func() {
		for {
			select {
			case <-ticker.C:
				log.Println("Doing stuff with " + config.FileName)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
	time.Sleep(30 * time.Second)
	close(quit)
}
