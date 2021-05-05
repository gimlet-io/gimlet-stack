package commands

import (
	"fmt"
	"github.com/gimlet-io/gimlet-stack/template"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path/filepath"
)

var GenerateCmd = cli.Command{
	Name:      "generate",
	Usage:     "Generates Kubernetes resources from stack.yaml",
	UsageText: `stack generate`,
	Action:    generate,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "config",
			Aliases:  []string{"c"},
		},
		&cli.StringFlag{
			Name:     "target-path",
			Aliases:  []string{"p"},
		},
	},
}

func generate(c *cli.Context) error {

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

	files, err := template.GenerateFromStackYaml(stackConfig)
	if err != nil {
		return fmt.Errorf("cannot parse stack config file: %s", err.Error())
	}

	targetPath := c.String("target-path")
	for path, content := range files {
		err := os.MkdirAll(filepath.Dir(path), 0775)
		if err != nil {
			return fmt.Errorf("cannot write stack: %s", err.Error())
		}

		err = ioutil.WriteFile(filepath.Join(targetPath, path), []byte(content), 0664)
		if err != nil {
			return fmt.Errorf("cannot write stack: %s", err.Error())
		}
	}

	return nil
}
