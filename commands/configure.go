package commands

import (
	"fmt"
	"github.com/gimlet-io/gimlet-stack/template"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

var ConfigureCmd = cli.Command{
	Name:      "configure",
	Usage:     "Configures Kubernetes resources and writes a stack.yaml",
	UsageText: `stack configure`,
	Action:    configure,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "config",
			Aliases:  []string{"c"},
		},
		&cli.StringFlag{
			Name:     "stack-repo",
			Aliases:  []string{"r"},
		},
	},
}

func configure(c *cli.Context) error {
	stackConfigPath := c.String("config")
	if stackConfigPath == "" {
		stackConfigPath = "stack.yaml"
	}
	stackConfigYaml, err := ioutil.ReadFile(stackConfigPath)
	if err != nil {
		return fmt.Errorf("cannot read stack config file: %s", err.Error())
	}

	var stackConfig template.StackConfig
	err = yaml.Unmarshal([]byte(stackConfigYaml), &stackConfig)
	if err != nil {
		return fmt.Errorf("cannot parse stack config file: %s", err.Error())
	}

	stackRepoURL := c.String("stack-repo")
	if stackRepoURL == "" {
		stackRepoURL = "https://github.com/gimlet-io/gimlet-stack-reference.git"
	}

	stackDefinitionYaml, err := template.StackDefinitionFromRepo(stackRepoURL)
	if err != nil {
		return fmt.Errorf("cannot get stack definition: %s", err.Error())
	}
	var stackDefinition template.StackDefinition
	err = yaml.Unmarshal([]byte(stackDefinitionYaml), &stackDefinition)
	if err != nil {
		return fmt.Errorf("cannot parse stack definition: %s", err.Error())
	}

	_, err = template.Configure(stackDefinition, stackConfig)
	if err != nil {
		return fmt.Errorf("cannot configure stack: %s", err.Error())
	}

	return nil
}
