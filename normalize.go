package main

import "github.com/envato/iamy/iamy"

type NormalizeCommandInput struct {
	Dir       string
	CanDelete bool
}

func NormalizeCommand(ui Ui, input NormalizeCommandInput) {
	if *dryRun {
		ui.Fatal("Dry-run mode not supported for normalize")
	}

	yaml := iamy.YamlLoadDumper{
		Dir: input.Dir,
	}

	allDataFromYaml, err := yaml.Load()
	if err != nil {
		ui.Fatal(err)
		return
	}

	for _, account := range allDataFromYaml {
		ui.Printf("Normalizing %s (%s)", account.Account.Alias, account.Account.Id)

		err = yaml.Dump(&account, input.CanDelete)
		if err != nil {
			ui.Error.Fatal(err)
		}
	}
}
