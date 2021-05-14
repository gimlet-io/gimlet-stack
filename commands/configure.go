package commands

import (
	"bytes"
	"fmt"
	"github.com/enescakir/emoji"
	"github.com/gimlet-io/gimlet-stack/template"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path/filepath"
)

var ConfigureCmd = cli.Command{
	Name:      "configure",
	Usage:     "Configures Kubernetes resources and writes a stack.yaml",
	UsageText: `stack configure`,
	Action:    configure,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
		},
		&cli.StringFlag{
			Name:    "stack-repo",
			Aliases: []string{"r"},
		},
	},
}

func configure(c *cli.Context) error {
	var stackConfig template.StackConfig
	stackConfigPath := c.String("config")

	if stackConfigPath == "" {
		stackConfigPath = "stack.yaml"
	}

	absolutePath, err := filepath.Abs(stackConfigPath)
	if err != nil {
		return fmt.Errorf("cannot parse stack config file: %s", err.Error())
	}

	if _, err := os.Stat(absolutePath); err == nil {
		stackConfigYaml, err := ioutil.ReadFile(stackConfigPath)
		if err != nil {
			return fmt.Errorf("cannot read stack config file: %s", err.Error())
		}
		err = yaml.Unmarshal(stackConfigYaml, &stackConfig)
		if err != nil {
			return fmt.Errorf("cannot parse stack config file: %s", err.Error())
		}
	}

	stackRepoURL := c.String("stack-repo")
	if stackRepoURL != "" {
		stackConfig.Stack.Repository = stackRepoURL
	}
	if stackConfig.Stack.Repository == "" {
		stackConfig.Stack.Repository = "https://github.com/gimlet-io/gimlet-stack-reference.git"
	}

	stackDefinitionYaml, err := template.StackDefinitionFromRepo(stackConfig.Stack.Repository)
	if err != nil {
		return fmt.Errorf("cannot get stack definition: %s", err.Error())
	}
	var stackDefinition template.StackDefinition
	err = yaml.Unmarshal([]byte(stackDefinitionYaml), &stackDefinition)
	if err != nil {
		return fmt.Errorf("cannot parse stack definition: %s", err.Error())
	}

	updatedStackConfig, err := template.Configure(stackDefinition, stackConfig)
	if err != nil {
		return fmt.Errorf("cannot configure stack: %s", err.Error())
	}

	updatedStackConfigBuffer := bytes.NewBufferString("")
	e := yaml.NewEncoder(updatedStackConfigBuffer)
	e.SetIndent(2)
	e.Encode(updatedStackConfig)

	updatedStackConfigString := "---\n" + updatedStackConfigBuffer.String()
	err = ioutil.WriteFile(stackConfigPath, []byte(updatedStackConfigString), 0666)
	if err != nil {
		return fmt.Errorf("cannot write stack file %s", err)
	}

	fmt.Println("---")
	fmt.Println(updatedStackConfigString)

	fmt.Fprintf(os.Stderr, "%v Written to %s \n\n", emoji.FileFolder, stackConfigPath)

	return nil
}
