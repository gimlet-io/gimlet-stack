package template

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"testing"
)

func Test_Configure(t *testing.T) {
	stackConfigYaml := `
stack:
  repository: "..."
config:
  nginx:
    enabled: true
`

	stackDefinitionYaml := `

`

	var stackConfig StackConfig
	err := yaml.Unmarshal([]byte(stackConfigYaml), &stackConfig)
	assert.Nil(t, err)

	var stackDefinition StackDefinition
	err = yaml.Unmarshal([]byte(stackDefinitionYaml), &stackDefinition)
	assert.Nil(t, err)

	_, err = Configure(stackDefinition, stackConfig)
}
