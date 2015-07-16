package main

import (
	"log"
	"time"

	"github.com/sparrc/stash/config"
)

func main() {

	done := make(chan bool)
	go func() {
		time.Sleep(10 * time.Second)
		done <- true
	}()

	configMngr := config.NewMngr()
	log.Println("Waiting " + configMngr.FileName)

	<-done
}
