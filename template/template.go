package template

import (
	"bytes"
	"github.com/Masterminds/sprig/v3"
	"strings"
	"text/template"
)

func Generate(
	stackTemplate map[string]string,
	values map[string]interface{},
) (map[string]string, error) {
	generatedFiles := map[string]string{}

	for path, fileContent := range stackTemplate {
		templates, err := template.New(path).Funcs(sprig.TxtFuncMap()).Parse(fileContent)
		if err != nil {
			return nil, err
		}

		var templated bytes.Buffer
		err = templates.Execute(&templated, values)
		if err != nil {
			return nil, err
		}

		// filter empty, or white space only files
		if len(strings.TrimSpace(templated.String())) != 0 {
			generatedFiles[path] = templated.String()
		}
	}

	return generatedFiles, nil
}
