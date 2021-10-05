package main

import (
	"fmt"

	"github.com/envato/iamy/iamy"
)

type PullCommandInput struct {
	Dir                  string
	CanDelete            bool
	HeuristicCfnMatching bool
	SkipTagged           []string
	IncludeTagged        []string
	SkipPaths            []string
}

func PullCommand(ui Ui, input PullCommandInput) {
	aws := iamy.AwsFetcher{
		Debug:                ui.Debug,
		HeuristicCfnMatching: input.HeuristicCfnMatching,
		SkipTagged:           input.SkipTagged,
		IncludeTagged:        input.IncludeTagged,
		SkipPaths:            input.SkipPaths,
	}
	data, err := aws.Fetch()
	if err != nil {
		ui.Error.Fatal(fmt.Printf("%s", err))
	}

	yaml := iamy.YamlLoadDumper{
		Dir: input.Dir,
	}
	err = yaml.Dump(data, input.CanDelete)
	if err != nil {
		ui.Error.Fatal(err)
	}
}
