package main

import (
	"log"

	"github.com/cameronsparr/stash/config"

	fsnotify "github.com/cameronsparr/stash/Godeps/_workspace/src/gopkg.in/fsnotify.v1"
)

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	configMngr := config.NewMngr()
	log.Println("Watching " + configMngr.FileName)

	err = watcher.Add(configMngr.FileName)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
