package template

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_FilterEmptyFiles(t *testing.T) {
	generatedFiles, err := Generate(map[string]string{"empty-file.yaml": ""}, map[string]interface{}{})
	assert.Nil(t, err)
	assert.Empty(t, generatedFiles)

	generatedFiles, err = Generate(map[string]string{"whitespace-only.yaml": "\n"}, map[string]interface{}{})
	assert.Nil(t, err)
	assert.Empty(t, generatedFiles)
}

func Test_BasicTemplating(t *testing.T) {
	stackTemplate := map[string]string{
		"template.yaml": `
{{- if .Enabled }}
---
apiVersion: helm.toolkit.fluxcd.io/v2beta1
kind: HelmRelease
{{- end }}
`,
	}

	generatedFiles, err := Generate(stackTemplate, map[string]interface{}{"Enabled": false})
	assert.Nil(t, err)
	assert.Empty(t, generatedFiles)

	generatedFiles, err = Generate(stackTemplate, map[string]interface{}{"Enabled": true})
	assert.Nil(t, err)
	assert.NotEmpty(t, generatedFiles)
}
