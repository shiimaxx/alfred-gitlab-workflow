package main

import (
	"io"

	"fmt"

	"flag"

	"os"

	"github.com/shiimaxx/alfred-gitlab-workflow/workflow"
)

// Exit codes are int values that represent an exit code for a particular error.
const (
	ExitCodeOK    int = 0
	ExitCodeError int = 1 + iota
)

// CLI is the command line object
type CLI struct {
	// outStream and errStream are the stdout and stderr
	// to write message from the CLI.
	outStream, errStream io.Writer
}

// Run invokes the CLI with the given arguments.
func (c *CLI) Run(args []string) int {
	var (
		setURL   bool
		setToken bool
	)
	flags := flag.NewFlagSet("alfred-gitlab-workflow", flag.ContinueOnError)
	flags.SetOutput(c.outStream)
	flags.BoolVar(&setURL, "set-url", false, "set endpoint url")
	flags.BoolVar(&setToken, "set-token", false, "set personal access token")
	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeError
	}

	if setURL {
		f, err := os.OpenFile("endpoint_url", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0664)
		if err != nil {
			fmt.Fprint(c.errStream, err)
			return ExitCodeError
		}
		url := flags.Args()[0]
		if _, err := f.Write([]byte(url)); err != nil {
			fmt.Fprint(c.errStream, err)
			return ExitCodeError
		}
		if err := f.Close(); err != nil {
			fmt.Fprint(c.errStream, err)
			return ExitCodeError
		}
		return ExitCodeOK
	}

	if setToken {
		return ExitCodeOK
	}

	fmt.Fprint(c.outStream, workflow.Run())
	return ExitCodeOK
}
