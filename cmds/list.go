package cmds

import (
	"fmt"
)

var List = &Command{
	Usage: "list [arguments]",
	Short: "list backup destinations & their associated folders.",
	Long: `
Adds a backup folder to an existing backup destination, or a backup destination.
`,
	Run: runList,
}

func runList(cmd *Command, args []string) {
	fmt.Println("Running List Command")
}
