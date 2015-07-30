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
	Usage: "destination [ add | list | delete | help ]",
	Short: "Add, list, or delete backup destinations.",
	Long: `
Usage:

	stash destination [command]

The commands are:

	add	Add a backup destination
	list	List configured backup destinations
	delete	Delete a backup destination & associated folders
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
	case "delete", "del":
		runDelete(args)
	case "help":
		cmd.LongUsageExit()
	default:
		color.Red("Invalid subcommand: %s", args[0])
		os.Exit(1)
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
		addGoogle()
	case "":
		return
	default:
		color.Red("Invalid option, exiting")
		os.Exit(1)
	}
}

func addAmazon() {
	config := stash.NewConfig()
	reader := bufio.NewReader(os.Stdin)
	confEntry := stash.ConfigEntry{
		Name:        userInputName(reader, config),
		Folders:     userInputFolders(reader),
		Type:        "Amazon",
		Credentials: userInputCredentials(reader),
		Frequency:   userInputFrequency(reader),
	}
	if err := config.AddDestination(confEntry); err != nil {
		color.Red("Fatal error adding backup destination: ", err)
		os.Exit(1)
	}
}

func addGoogle() {
	return
}

func userInputName(reader *bufio.Reader, confFile *stash.Config) string {
	fmt.Print("Specify a name for this destination: ")
	text, _ := reader.ReadString('\n')
	name := strings.TrimSpace(text)
	if confFile.IsDuplicateEntry(stash.ConfigEntry{Name: name}) {
		color.Red("Attempted to add duplicate entry [%s], "+
			"if you were trying to add folders to an existing backup "+
			"destination, use 'stash folder add'", name)
		os.Exit(1)
	}
	return name
}

func userInputFolders(reader *bufio.Reader) []string {
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
	if _, err := os.Stat(dir); err != nil {
		red := color.New(color.FgRed).SprintfFunc()
		var e string
		if os.IsNotExist(err) {
			e = red("Directory %s does not exist", dir)
		} else {
			e = red("Error accessing %s, do you have read permission?", dir)
		}
		return false, errors.New(e)
	}
	return true, nil
}

func userInputCredentials(reader *bufio.Reader) map[string]string {
	return map[string]string{"key": "supersecret"}
}

func userInputFrequency(reader *bufio.Reader) time.Duration {
	color.Blue("Backup Frequency, use a short string like: 30m or 2h43m10s")
	color.Blue("Valid time units are: h, m, s")
	fmt.Println("")
	fmt.Print("Backup Every: ")
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	d, err := time.ParseDuration(text)
	if err != nil {
		color.Red("Error parsing frequency, try again")
		return userInputFrequency(reader)
	}
	return d
}

func runList(args []string) {
	config := stash.NewConfig()
	col := color.New(color.FgMagenta)
	color.New(color.FgBlue, color.Bold).Println("Current Backup Destinations:")
	fmt.Println()
	for _, entry := range config.Entries {
		col.Printf("Name:		")
		fmt.Println(entry.Name)
		col.Printf("Folders:	")
		fmt.Println(entry.Folders)
		col.Printf("Type:		")
		fmt.Println(entry.Type)
		col.Printf("Frequency:	")
		fmt.Println("Every ", entry.Frequency)
		col.Printf("Last Backup:	")
		// .Format takes an example string for how the stamp should look
		fmt.Println(entry.LastBak.Format("Jan 2, 2006 at 3:04pm (MST)"))
		fmt.Println()
	}
}

func runDelete(args []string) {
	config := stash.NewConfig()
	reader := bufio.NewReader(os.Stdin)
	color.Red("WARNING! Deleted entries cannot be recovered")
	color.Blue("Please choose one of the following entries: ")
	for _, entry := range config.Entries {
		color.Magenta("	%s", entry.Name)
	}
	fmt.Println()
	fmt.Print("Entry to delete (case-sensitive): ")
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	err := config.DeleteEntry(text)
	if err != nil {
		color.Red("Fatal error deleting an entry: %s", err)
		os.Exit(1)
	}
}
