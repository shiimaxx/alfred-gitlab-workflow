package main

import (
	"io"

	"fmt"

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
	fmt.Fprint(c.outStream, workflow.Run())
	return ExitCodeOK
}
