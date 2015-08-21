package main

import (
	"flag"
	"strings"

	"github.com/99designs/iamy/Godeps/_workspace/src/github.com/mitchellh/cli"
	"github.com/99designs/iamy/iamy"
)

type DumpCommand struct {
	Ui cli.Ui
}

func (c *DumpCommand) Run(args []string) int {
	var dir string
	flagSet := flag.NewFlagSet("dump", flag.ContinueOnError)
	flagSet.StringVar(&dir, "dir", "", "Directory to write files to")
	flagSet.Usage = func() { c.Ui.Output(c.Help()) }

	if err := flagSet.Parse(args); err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	dir, err := getDirOrDefault(dir)
	if err != nil {
		c.Ui.Error(err.Error())
		return 2
	}

	// fetch data from AWS
	data, err := iamy.Aws.Fetch()
	if err != nil {
		c.Ui.Error(err.Error())
		return 3
	}

	// dump data to dir
	iamy.Yaml.Dir = dir
	err = iamy.Yaml.Dump(data)
	if err != nil {
		c.Ui.Error(err.Error())
		return 4
	}

	return 0
}

func (c *DumpCommand) Help() string {
	helpText := `
Usage: iamy dump [-dir <output dir>]
  Dumps users, groups and policies to files
`
	return strings.TrimSpace(helpText)
}

func (c *DumpCommand) Synopsis() string {
	return "Dumps users, groups and policies to files"
}
