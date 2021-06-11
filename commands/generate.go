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

var GenerateCmd = cli.Command{
	Name:      "generate",
	Usage:     "Generates Kubernetes resources from stack.yaml",
	UsageText: `stack generate`,
	Action:    generate,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
		},
		&cli.StringFlag{
			Name:    "target-path",
			Aliases: []string{"p"},
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
	err = yaml.Unmarshal(stackConfigYaml, &stackConfig)
	if err != nil {
		return fmt.Errorf("cannot parse stack config file: %s", err.Error())
	}

	err = lockVersionIfNotLocked(stackConfig, stackConfigPath)
	if err != nil {
		return fmt.Errorf("couldn't lock stack version: %s", err.Error())
	}

	checkForUpdates(stackConfig)

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

	fmt.Fprintf(os.Stderr, "%v  Generated\n", emoji.CheckMark)

	return nil
}

func checkForUpdates(stackConfig template.StackConfig) {
	currentTagString := template.CurrentVersion(stackConfig.Stack.Repository)
	if currentTagString != "" {
		versionsSince, err := template.VersionsSince(stackConfig.Stack.Repository, currentTagString)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\n%v  Cannot check for updates \n\n", emoji.Warning)
		}

		if len(versionsSince) > 0 {
			fmt.Fprintf(os.Stderr, "\n%v  Stack update available. Run `stack update --check` for details. \n\n", emoji.Warning)
		}
	}
}

func lockVersionIfNotLocked(stackConfig template.StackConfig, stackConfigPath string) error {
	locked, err := template.IsVersionLocked(stackConfig)
	if err != nil {
		return fmt.Errorf("cannot check version: %s", err.Error())
	}
	if !locked {
		latestTag, _ := template.LatestVersion(stackConfig.Stack.Repository)
		if latestTag != "" {
			stackConfig.Stack.Repository = stackConfig.Stack.Repository + "?tag=" + latestTag

			updatedStackConfigBuffer := bytes.NewBufferString("")
			e := yaml.NewEncoder(updatedStackConfigBuffer)
			e.SetIndent(2)
			e.Encode(stackConfig)

			updatedStackConfigString := "---\n" + updatedStackConfigBuffer.String()
			err = ioutil.WriteFile(stackConfigPath, []byte(updatedStackConfigString), 0666)
			if err != nil {
				return fmt.Errorf("cannot write stack file %s", err)
			}

			fmt.Fprintf(os.Stderr, "%v  Stack version is locked to %s \n\n", emoji.Warning, latestTag)
		}
	}

	return nil
}
