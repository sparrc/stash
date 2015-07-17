package main

import (
	"flag"
	"io"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/fatih/color"

	"github.com/sparrc/stash/subcmd"
)

// Commands lists the available commands and help topics.
// The order here is the order in which they are printed
// by 'stash help'.
var commands = []*subcmd.Command{
	subcmd.Destination,
	subcmd.Folder,
	subcmd.Daemon,
}

func main() {
	flag.Usage = usageExit
	flag.Parse()
	log.SetFlags(0)
	log.SetPrefix("DEBUG: ")
	args := flag.Args()
	if len(args) < 1 {
		usageExit()
	}

	if args[0] == "help" {
		help(args[1:])
		return
	}

	for _, cmd := range commands {
		if cmd.Name() == args[0] {
			cmd.Flag.Usage = func() { cmd.UsageExit() }
			cmd.Flag.Parse(args[1:])
			cmd.Run(cmd, cmd.Flag.Args())
			return
		}
	}

	color.Red("stash: unknown command %q\n", args[0])
	color.Red("Run 'stash help' for usage.\n")
	os.Exit(2)
}

var usageTemplate = `
Stash is a CLI tool for managing cloud backups to Amazon AWS and Google Cloud.

Usage:

	stash command [arguments]

The commands are:
{{range .}}
	{{.Name | printf "%-8s"}}	{{.Short}}{{end}}

Use "stash help [command]" for more information about a command.
`

var helpTemplate = `
{{.Long | trim}}
`

func help(args []string) {
	if len(args) == 0 {
		printUsage(os.Stdout)
		return
	}
	if len(args) != 1 {
		color.Red("usage: stash help command\n\n")
		color.Red("Too many arguments given.\n")
		os.Exit(2)
	}
	for _, cmd := range commands {
		if cmd.Name() == args[0] {
			tmpl(os.Stdout, helpTemplate, cmd)
			return
		}
	}
}

func usageExit() {
	printUsage(os.Stderr)
	os.Exit(2)
}

func printUsage(w io.Writer) {
	tmpl(w, usageTemplate, commands)
}

// tmpl executes the given template text on data, writing the result to w.
func tmpl(w io.Writer, text string, data interface{}) {
	t := template.New("top")
	t.Funcs(template.FuncMap{
		"trim": strings.TrimSpace,
	})
	template.Must(t.Parse(strings.TrimSpace(text) + "\n\n"))
	if err := t.Execute(w, data); err != nil {
		panic(err)
	}
}
