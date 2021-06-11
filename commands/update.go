package commands

import (
	"bytes"
	"fmt"
	markdown "github.com/MichaelMure/go-term-markdown"
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
			Name:    "config",
			Aliases: []string{"c"},
		},
		&cli.BoolFlag{
			Name: "check",
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

	check := c.Bool("check")

	currentTagString := template.CurrentVersion(stackConfig.Stack.Repository)
	versionsSince, err := template.VersionsSince(stackConfig.Stack.Repository, currentTagString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n%v  Cannot check for updates \n\n", emoji.Warning)
	}

	if len(versionsSince) == 0 {
		fmt.Fprintf(os.Stderr, "\n%v  Already up to date \n\n", emoji.CheckMark)
		return nil
	}

	if check {
		fmt.Fprintf(os.Stderr, "%v  New version available: \n\n", emoji.Books)
		err := printChangeLog(stackConfig, versionsSince)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\n%v %s \n\n", emoji.Warning, err)
		}
		fmt.Fprintf(os.Stderr, "\n")
	} else {
		latestTag, _ := template.LatestVersion(stackConfig.Stack.Repository)
		if latestTag != "" {
			fmt.Fprintf(os.Stderr, "%v  Stack version is updating to %s... \n\n", emoji.HourglassNotDone, latestTag)
			stackConfig.Stack.Repository = template.RepoUrlWithoutVersion(stackConfig.Stack.Repository) + "?tag=" + latestTag
			err = writeStackConfig(stackConfig, stackConfigPath)
			if err != nil {
				return fmt.Errorf("cannot write stack file %s", err)
			}
			fmt.Fprintf(os.Stderr, "%v   Config updated. \n\n", emoji.CheckMark)
			fmt.Fprintf(os.Stderr, "%v   Run `stack generate` to render resources with the updated stack. \n\n", emoji.Warning)
			fmt.Fprintf(os.Stderr, "%v  Change log:\n\n", emoji.Books)
			err = printChangeLog(stackConfig, versionsSince)
			if err != nil {
				fmt.Fprintf(os.Stderr, "\n%v %s \n\n", emoji.Warning, err)
			}
			fmt.Fprintf(os.Stderr, "\n")
		} else {
			fmt.Fprintf(os.Stderr, "%v  cannot find latest version\n", emoji.CrossMark)
		}
	}

	return nil
}

func writeStackConfig(stackConfig template.StackConfig, stackConfigPath string) error {
	updatedStackConfigBuffer := bytes.NewBufferString("")
	e := yaml.NewEncoder(updatedStackConfigBuffer)
	e.SetIndent(2)
	e.Encode(stackConfig)

	updatedStackConfigString := "---\n" + updatedStackConfigBuffer.String()
	return ioutil.WriteFile(stackConfigPath, []byte(updatedStackConfigString), 0666)
}

func printChangeLog(stackConfig template.StackConfig, versions []string) error {
	for _, version := range versions {
		fmt.Fprintf(os.Stderr, "   - %s \n", version)

		repoUrl := stackConfig.Stack.Repository
		repoUrl = template.RepoUrlWithoutVersion(repoUrl)
		repoUrl = repoUrl + "?tag=" + version

		stackDefinitionYaml, err := template.StackDefinitionFromRepo(repoUrl)
		if err != nil {
			return fmt.Errorf("cannot get stack definition: %s", err.Error())
		}
		var stackDefinition template.StackDefinition
		err = yaml.Unmarshal([]byte(stackDefinitionYaml), &stackDefinition)
		if err != nil {
			return fmt.Errorf("cannot parse stack definition: %s", err.Error())
		}

		if stackDefinition.ChangLog != "" {
			changeLog := markdown.Render(stackDefinition.ChangLog, 80, 6)
			fmt.Fprintf(os.Stderr, "%s\n", changeLog)
		}
	}

	return nil
}
