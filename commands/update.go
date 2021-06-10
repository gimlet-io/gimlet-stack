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
)

var UpdateCmd = cli.Command{
	Name:      "update",
	Usage:     "Updates the stack version in stack.yaml",
	UsageText: `stack update`,
	Action:    update,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "config",
			Aliases:  []string{"c"},
		},
	},
}

func update(c *cli.Context) error {
	stackConfigPath := c.String("config")
	if stackConfigPath == "" {
		stackConfigPath = "stack.yaml"
	}
	stackConfigYaml, err := ioutil.ReadFile(stackConfigPath)
	if err != nil {
		return fmt.Errorf("cannot read stack config file: %s", err.Error())
	}

	var stackConfig template.StackConfig
	err = yaml.Unmarshal(stackConfigYaml, &stackConfig)
	if err != nil {
		return fmt.Errorf("cannot parse stack config file: %s", err.Error())
	}

	latestTag, _ := template.LatestVersion(stackConfig.Stack.Repository)
	if latestTag != "" {
		stackConfig.Stack.Repository = template.RepoUrlWithoutVersion(stackConfig.Stack.Repository) + "?tag=" + latestTag

		updatedStackConfigBuffer := bytes.NewBufferString("")
		e := yaml.NewEncoder(updatedStackConfigBuffer)
		e.SetIndent(2)
		e.Encode(stackConfig)

		updatedStackConfigString := "---\n" + updatedStackConfigBuffer.String()
		err = ioutil.WriteFile(stackConfigPath, []byte(updatedStackConfigString), 0666)
		if err != nil {
			return fmt.Errorf("cannot write stack file %s", err)
		}

		fmt.Fprintf(os.Stderr, "%v  Stack version is updated to %s \n\n", emoji.CheckMark, latestTag)
	} else {
		fmt.Fprintf(os.Stderr, "%v  cannot find latest version\n", emoji.CrossMark)
	}

	return nil
}
