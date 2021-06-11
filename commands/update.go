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
		fmt.Fprintf(os.Stderr, "%v  Already up to date \n\n", emoji.CheckMark)
		return nil
	}

	if check {
		fmt.Fprintf(os.Stderr, "%v  New version available: \n\n", emoji.Books)
		printChangeLog(versionsSince)
		fmt.Fprintf(os.Stderr, "\n")
	} else {
		latestTag, _ := template.LatestVersion(stackConfig.Stack.Repository)
		if latestTag != "" {
			fmt.Fprintf(os.Stderr, "%v  Stack version is updating to %s... \n\n", emoji.HourglassNotDone, latestTag)
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

			fmt.Fprintf(os.Stderr, "%v   Updated.\n\n", emoji.CheckMark)
			fmt.Fprintf(os.Stderr, "%v   Run `stack generate` to render resources with the updated stack. \n\n", emoji.Warning)
			fmt.Fprintf(os.Stderr, "%v  Change log:\n\n", emoji.Books)
			printChangeLog(versionsSince)
			fmt.Fprintf(os.Stderr, "\n")
		} else {
			fmt.Fprintf(os.Stderr, "%v  cannot find latest version\n", emoji.CrossMark)
		}
	}

	return nil
}

func printChangeLog(versions []string) {
	for _, version := range versions {
		fmt.Fprintf(os.Stderr, "   - %s \n", version)
	}
}
