package cmds

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/cameronsparr/stash/config"
)

// Destination specifies the 'stash destination' command.
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
	text = strings.TrimSpace(text)
	switch text {
	case "1":
		addAmazon()
	case "2":
		return
	}
}

func addAmazon() {
	// TODO: implement
	conf := config.NewConfigMngr()
	reader := bufio.NewReader(os.Stdin)
	confEntry := config.Entry{
		Name:        getName(reader),
		Folders:     getFolders(reader),
		Type:        "Amazon",
		Credentials: getCredentials(reader),
	}
	conf.AddDestination(confEntry)
}

func getName(reader *bufio.Reader) string {
	// TODO: check if user is trying to add a duplicate destination
	fmt.Print("Specify a name for this destination: ")
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func getFolders(reader *bufio.Reader) []string {
	// TODO: check that folders are valid
	fmt.Print("Specify directories to backup (space-separated): ")
	text, _ := reader.ReadString('\n')
	return strings.Split(strings.TrimSpace(text), " ")
}

func getCredentials(reader *bufio.Reader) map[string]string {
	return map[string]string{"key": "supersecret"}
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
