package cmds

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

var Add = &Command{
	Usage: "add [destination | folder] [arguments]",
	Short: "add a backup destination or a folder to an existing backup destination",
	Long: `
Adds a backup folder to an existing backup destination, or a backup destination.
`,
	Run: runAdd,
}

func runAdd(cmd *Command, args []string) {
	if len(args) == 0 {
		log.Println("add requires one of [destination | folder]")
		cmd.UsageExit()
	} else if args[0] == "destination" {
		runAddDestination(args)
	} else if args[0] == "folder" {
		runAddFolder(args)
	}

}

func runAddDestination(args []string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(`
Which type of backup destination would you like to add?
	1. Amazon
	2. Google Cloud
`)
	text, _ := reader.ReadString('\n')
	fmt.Println(text)
}

func runAddFolder(args []string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(`
This is not implemented yet, but do you love marutaro? [Y/y]
`)
	text, _ := reader.ReadString('\n')
	fmt.Println(text)
}
