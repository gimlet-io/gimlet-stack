package commands

import (
	"github.com/urfave/cli/v2"
)

var GenerateCmd = cli.Command{
	Name:      "generate",
	Usage:     "Generates Kubernetes resources from stack.yaml",
	UsageText: `stack generate`,
	Action:    generate,
	Flags: []cli.Flag{
	},
}

func generate(c *cli.Context) error {
	return nil
}
