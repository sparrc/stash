package cmds

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/cameronsparr/stash/config"
)

var Destination = &Command{
	Usage: "destination [add | list | remove] [arguments]",
	Short: "Add, list, or remove backup destinations.",
	Long: `
Add, list, or remove backup destinations.
`,
	Run: runDestination,
}

func runDestination(cmd *Command, args []string) {
	if len(args) == 0 {
		cmd.UsageExit()
	} else if args[0] == "add" {
		runAdd(args)
	} else if args[0] == "list" {
		runList(args)
	} else if args[0] == "remove" {
		runRemove(args)
	}
}

func runAdd(args []string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Which type of backup destination would you like to add?")
	fmt.Println("	1. Amazon (S3 or Glacier)")
	fmt.Println("	2. Google Cloud")
	fmt.Println("")
	fmt.Print("Choose an option [1-2]: ")
	text, _ := reader.ReadString('\n')
	log.Println(text)

	// TODO: Actually do the configuration
	conf := config.NewConfig()
	awsDest := config.AmazonDestination{DestinationName: "AWS"}
	conf.AddDestination(&awsDest)
	conf.LoadConfigFile()
}

func runList(args []string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("This is not implemented yet, but do you love marutaro? [Y/y]")
	text, _ := reader.ReadString('\n')
	fmt.Println(text)
}

func runRemove(args []string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("This is not implemented yet, but do you love marutaro? [Y/y]")
	text, _ := reader.ReadString('\n')
	fmt.Println(text)
}
