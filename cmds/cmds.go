package cmds

import (
    "fmt"
    "flag"
    "strings"
    "os"
)

// A Command is an implementation of a stash command
// like stash add or stash list.
type Command struct {
    // Run runs the command.
    // The args are the arguments after the command name.
    Run func(cmd *Command, args []string)

    // Usage is the one-line usage message.
    // The first word in the line is taken to be the command name.
    Usage string

    // Short is the short description shown in the 'stash help' output.
    Short string

    // Long is the long message shown in the
    // 'stash help <this-command>' output.
    Long string

    // Flag is a set of flags specific to this command.
    Flag flag.FlagSet
}

func (c *Command) Name() string {
    name := c.Usage
    i := strings.Index(name, " ")
    if i >= 0 {
        name = name[:i]
    }
    return name
}

func (c *Command) UsageExit() {
    fmt.Fprintf(os.Stderr, "Usage: stash %s\n\n", c.Usage)
    fmt.Fprintf(os.Stderr, "Run 'stash help %s' for help.\n", c.Name())
    os.Exit(2)
}