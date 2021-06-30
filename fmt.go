package main

import "github.com/envato/iamy/iamy"

type FormatCommandInput struct {
	Dir       string
	CanDelete bool
}

func FormatCommand(ui Ui, input FormatCommandInput) {
	if *dryRun {
		ui.Fatal("Dry-run mode not supported for fmt")
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
		ui.Printf("Formatting %s (%s)", account.Account.Alias, account.Account.Id)

		err = yaml.Dump(&account, input.CanDelete)
		if err != nil {
			ui.Error.Fatal(err)
		}
	}
}
