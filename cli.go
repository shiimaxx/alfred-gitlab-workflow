package main

import (
	"io"
	"fmt"
	"flag"
	"os"
	"io/ioutil"

	"github.com/shiimaxx/alfred-gitlab-workflow/workflow"
	"github.com/keybase/go-keychain"
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
		item := keychain.NewGenericPassword("alfred-gitlab-workflow", "", "", []byte(flags.Args()[0]), "")
		item.SetAccessible(keychain.AccessibleWhenUnlocked)
		if err := keychain.AddItem(item); err != nil {
			fmt.Fprint(c.errStream, err)
			return ExitCodeError
		}
		return ExitCodeOK
	}

	var url string
	if _, err := os.Stat("endpoint_url"); ! os.IsNotExist(err) {
		d, err := ioutil.ReadFile("endpoint_url")
		if err != nil {
			fmt.Fprint(c.errStream, err)
			return ExitCodeError
		}
		url = string(d)
	}

	b, err := keychain.GetGenericPassword("alfred-gitlab-workflow", "", "", "")
	if err != nil {
		fmt.Fprint(c.errStream, err)
		return ExitCodeError
	}

	fmt.Fprint(c.outStream, workflow.Run(url, string(b)))
	return ExitCodeOK
}
