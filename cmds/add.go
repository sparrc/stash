package cmds

import (
    "fmt"
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
    fmt.Println("Running Add Command")
}