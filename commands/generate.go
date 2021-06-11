package commands

import (
	"bytes"
	"fmt"
	markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/enescakir/emoji"
	"github.com/epiclabs-io/diff3"
	"github.com/gimlet-io/gimlet-stack/template"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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
	stackConfig, err := readStackConfig(stackConfigPath)
	if err != nil {
		return err
	}

	err = lockVersionIfNotLocked(stackConfig, stackConfigPath)
	if err != nil {
		return fmt.Errorf("couldn't lock stack version: %s", err.Error())
	}
	checkForUpdates(stackConfig)

	generatedFiles, err := template.GenerateFromStackYaml(stackConfig)
	if err != nil {
		return fmt.Errorf("cannot generate stack: %s", err.Error())
	}

	oldStackConfigPath := filepath.Join(filepath.Dir(stackConfigPath), ".stack", "old")
	oldStackConfig, err := readStackConfig(oldStackConfigPath)
	if err != nil {
		oldStackConfig = stackConfig
	}
	previousGenerationFiles, err := template.GenerateFromStackYaml(oldStackConfig)
	if err != nil {
		return fmt.Errorf("cannot generate stack: %s", err.Error())
	}

	targetPath := c.String("target-path")
	err = writeFilesAndPreserveCustomChanges(
		previousGenerationFiles,
		generatedFiles,
		targetPath,
	)
	if err != nil {
		fmt.Errorf("cannot write stack: %s", err.Error())
	}

	err = keepStackConfigUsedForGeneration(stackConfigPath, stackConfig)
	if err != nil {
		return fmt.Errorf("cannot write old stack config: %s", err.Error())
	}

	fmt.Fprintf(os.Stderr, "\n%v  Generated\n\n", emoji.CheckMark)

	stackDefinitionYaml, err := template.StackDefinitionFromRepo(stackConfig.Stack.Repository)
	if err != nil {
		return fmt.Errorf("cannot get stack definition: %s", err.Error())
	}
	var stackDefinition template.StackDefinition
	err = yaml.Unmarshal([]byte(stackDefinitionYaml), &stackDefinition)
	if err != nil {
		return fmt.Errorf("cannot parse stack definition: %s", err.Error())
	}

	if stackDefinition.ChangLog != "" {
		message := markdown.Render(stackDefinition.Message, 80, 6)
		fmt.Fprintf(os.Stderr, "%s\n", message)
	}

	return nil
}

func readStackConfig(stackConfigPath string) (template.StackConfig, error) {
	stackConfigYaml, err := ioutil.ReadFile(stackConfigPath)
	if err != nil {
		return template.StackConfig{}, fmt.Errorf("cannot read stack config file: %s", err.Error())
	}

	var stackConfig template.StackConfig
	err = yaml.Unmarshal(stackConfigYaml, &stackConfig)
	if err != nil {
		return template.StackConfig{}, fmt.Errorf("cannot parse stack config file: %s", err.Error())
	}
	return stackConfig, nil
}

func writeFilesAndPreserveCustomChanges(
	previousGenerationFiles map[string]string,
	generatedFiles map[string]string,
	targetPath string,
) error {
	for path, updated := range generatedFiles {
		err := os.MkdirAll(filepath.Dir(path), 0775)
		if err != nil {
			return fmt.Errorf("cannot write stack: %s", err.Error())
		}

		path = filepath.Join(targetPath, path)

		var existingContent string
		if _, err := os.Stat(path); err == nil {
			existingContentBytes, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("cannot read file %s: %s", path, err.Error())
			}
			existingContent = string(existingContentBytes)
		}

		var baseline string
		if val, ok := previousGenerationFiles[path]; ok {
			baseline = val
		}

		merged, err := diff3.Merge(strings.NewReader(existingContent), strings.NewReader(baseline), strings.NewReader(updated), true, "Your custom settings", "From stack generate")
		if err != nil {
			return fmt.Errorf("cannot merge %s: %s", path, err.Error())
		}
		mergedBuffer := new(strings.Builder)
		_, err = io.Copy(mergedBuffer, merged.Result)
		if err != nil {
			return fmt.Errorf("cannot merge %s: %s", path, err.Error())
		}

		mergedString := mergedBuffer.String()
		if !strings.HasSuffix(mergedString, "\n") {
			mergedString = mergedString + "\n"
		}

		err = ioutil.WriteFile(path, []byte(mergedString), 0664)
		if err != nil {
			return fmt.Errorf("cannot write stack: %s", err.Error())
		}
	}

	return nil
}

func keepStackConfigUsedForGeneration(
	stackConfigPath string,
	stackConfig template.StackConfig,
) error {
	stackBackupPath := filepath.Join(filepath.Dir(stackConfigPath), ".stack", "old")
	err := os.MkdirAll(filepath.Dir(stackBackupPath), 0775)
	if err != nil {
		return err
	}
	return writeStackConfig(stackConfig, stackBackupPath)
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
