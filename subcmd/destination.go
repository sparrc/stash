package subcmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"

	"github.com/sparrc/stash"
)

// Destination specifies the 'stash destination' command.
var Destination = &Command{
	Usage: "destination [ add | list | remove ]",
	Short: "Add, list, or remove backup destinations.",
	Long: `
Usage:

	stash destination [command]

The commands are:

	add	Add a backup destination
	list	List configured backup destinations
	remove	Remove a backup destination & associated folders
`,
	Run: runDestination,
}

func runDestination(cmd *Command, args []string) {
	if len(args) == 0 {
		cmd.UsageExit()
	}
	switch args[0] {
	case "add":
		runAdd(args)
	case "list":
		runList(args)
	case "remove":
		runRemove(args)
	}
}

func runAdd(args []string) {
	reader := bufio.NewReader(os.Stdin)
	color.Blue("Which type of backup destination would you like to add?")
	color.Blue("	1. Amazon (S3 or Glacier)")
	color.Blue("	2. Google Cloud")
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
	confFile := stash.NewConfig()
	reader := bufio.NewReader(os.Stdin)
	confEntry := stash.ConfigEntry{
		Name:        getName(reader, confFile),
		Folders:     getFolders(reader),
		Type:        "Amazon",
		Credentials: getCredentials(reader),
		Frequency:   getFrequency(reader),
	}
	confFile.AddDestination(confEntry)
}

func getName(reader *bufio.Reader, confFile *stash.Config) string {
	fmt.Print("Specify a name for this destination: ")
	text, _ := reader.ReadString('\n')
	name := strings.TrimSpace(text)
	if confFile.IsDuplicateEntry(stash.ConfigEntry{Name: name}) {
		color.Red("Attempted to add duplicate entry [%s], "+
			"if you were trying to add folders to an existing backup destination, "+
			"use 'stash folder add'", name)
		os.Exit(1)
	}
	return name
}

func getFolders(reader *bufio.Reader) []string {
	fmt.Print("Specify directories to backup (space-separated): ")
	text, _ := reader.ReadString('\n')
	dirs := strings.Split(strings.TrimSpace(text), " ")
	// Check that folders are valid and accessible
	for _, dir := range dirs {
		if _, err := isValidDirectory(dir); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	return dirs
}

func isValidDirectory(dir string) (bool, error) {
	red := color.New(color.FgRed).SprintfFunc()
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			e := red("Directory %s does not exist", dir)
			return false, errors.New(e)
		} else {
			e := red("Error accessing %s, do you have read permission?", dir)
			return false, errors.New(e)
		}
	}
	return true, nil
}

func getCredentials(reader *bufio.Reader) map[string]string {
	return map[string]string{"key": "supersecret"}
}

func getFrequency(reader *bufio.Reader) time.Duration {
	color.Blue("Backup Frequency, input as short string (ie, 30m or 2h43m10s)")
	color.Blue("Valid time units are s, m, h")
	fmt.Print("Backup Every: ")
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	d, err := time.ParseDuration(text)
	if err != nil {
		color.Red("Error parsing frequency, try again")
		return getFrequency(reader)
	}
	return d
}

func runList(args []string) {
	reader := bufio.NewReader(os.Stdin)
	color.Cyan("This is not implemented yet, but do you love marutaro? [Y/y]")
	text, _ := reader.ReadString('\n')
	fmt.Println(text)
}

func runRemove(args []string) {
	reader := bufio.NewReader(os.Stdin)
	color.Cyan("This is not implemented yet, but do you love marutaro? [Y/y]")
	text, _ := reader.ReadString('\n')
	fmt.Println(text)
}
