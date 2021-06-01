package main

import (
	"fmt"
	"github.com/gimlet-io/gimlet-stack/commands"
	"github.com/urfave/cli/v2"
	"os"

	"github.com/enescakir/emoji"
	"github.com/gimlet-io/gimlet-stack/version"
)

func main() {
	app := &cli.App{
		Name:                 "stack",
		Version:              version.String(),
		Usage:                "bootstrap curated Kubernetes stacks",
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			&commands.BootstrapCmd,
			&commands.GenerateCmd,
			&commands.ConfigureCmd,
			&commands.LintCmd,
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s %s\n", emoji.CrossMark, err.Error())
		os.Exit(1)
	}
}
