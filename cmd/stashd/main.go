package main

import (
	"log"
	"time"

	"github.com/sparrc/stash"
)

func main() {

	done := make(chan bool)
	go func() {
		time.Sleep(10 * time.Second)
		done <- true
	}()

	config := stash.NewConfig()
	log.Println("Waiting " + config.FileName)

	<-done
}
