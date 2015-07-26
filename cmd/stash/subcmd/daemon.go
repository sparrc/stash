package subcmd

import (
//	"fmt"
)

// Daemon specifies the 'stash daemon' command.
var Daemon = &Command{
	Usage: "daemon [arguments]",
	Short: "control the stash daemon process",
	Long: `
Control the stash daemon process.
`,
	Run: runDaemon,
}

func runDaemon(cmd *Command, args []string) {
	cmd.UsageExit()
}
