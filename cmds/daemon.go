package cmds

import (
	"fmt"
)

var Daemon = &Command{
	Usage: "daemon [arguments]",
	Short: "control the stash daemon process",
	Long: `
Control the stash daemon process.
`,
	Run: runDaemon,
}

func runDaemon(cmd *Command, args []string) {
	fmt.Println("Running Daemon Command")
}
