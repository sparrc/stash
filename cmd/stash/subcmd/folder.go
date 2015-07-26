package subcmd

import (
//	"fmt"
)

// Folder specified the 'stash folder' command.
var Folder = &Command{
	Usage: "folder [arguments]",
	Short: "folder backup destinations & their associated folders.",
	Long: `
Adds a backup folder to an existing backup destination, or a backup destination.
`,
	Run: runFolder,
}

func runFolder(cmd *Command, args []string) {
	cmd.UsageExit()
}
