package template

import (
	"bytes"
	"fmt"
	"github.com/Masterminds/sprig/v3"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-billy/v5/util"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
	giturl "github.com/whilp/git-urls"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type StackRef struct {
	Repository string `yaml:"repository" json:"repository"`
}

type StackConfig struct {
	Stack      StackRef               `yaml:"stack" json:"stack"`
	Config     map[string]interface{} `yaml:"config" json:"config"`
}

func GenerateFromStackYaml(stackConfig StackConfig) (map[string]string, error) {
	stackTemplates, err := cloneStackFromRepo(stackConfig.Stack.Repository)
	if err != nil {
		return nil, err
	}

	return generate(stackTemplates, stackConfig.Config)
}

func generate(
	stackTemplate map[string]string,
	values map[string]interface{},
) (map[string]string, error) {
	generatedFiles := map[string]string{}

	for path, fileContent := range stackTemplate {
		if path == "stack-definition.yaml" {
			continue
		}
		templates, err := template.New(path).Funcs(sprig.TxtFuncMap()).Parse(fileContent)
		if err != nil {
			return nil, err
		}

		var templated bytes.Buffer
		err = templates.Execute(&templated, values)
		if err != nil {
			return nil, err
		}

		// filter empty and white space only files
		if len(strings.TrimSpace(templated.String())) != 0 {
			generatedFiles[path] = templated.String()
		}
	}

	return generatedFiles, nil
}

func cloneStackFromRepo(repoURL string) (map[string]string, error) {
	gitAddress, err := giturl.ParseScp(repoURL)
	if err != nil {
		_, err2 := os.Stat(repoURL)
		if err2 != nil {
			return nil, fmt.Errorf("cannot parse stacks's git address: %s", err)
		} else {
			return loadStackFromFS(repoURL)
		}
	}
	gitUrl := strings.ReplaceAll(repoURL, gitAddress.RawQuery, "")
	gitUrl = strings.ReplaceAll(gitUrl, "?", "")

	fs := memfs.New()
	opts := &git.CloneOptions{
		URL: gitUrl,
	}
	repo, err := git.Clone(memory.NewStorage(), fs, opts)
	if err != nil {
		return nil, fmt.Errorf("cannot clone: %s", err)
	}
	worktree, err := repo.Worktree()
	if err != nil {
		return nil, fmt.Errorf("cannot get worktree: %s", err)
	}

	params, _ := url.ParseQuery(gitAddress.RawQuery)
	if v, found := params["sha"]; found {
		err = worktree.Checkout(&git.CheckoutOptions{
			Hash: plumbing.NewHash(v[0]),
		})
		if err != nil {
			return nil, fmt.Errorf("cannot checkout sha: %s", err)
		}
	}
	if v, found := params["tag"]; found {
		err = worktree.Checkout(&git.CheckoutOptions{
			Branch: plumbing.NewTagReferenceName(v[0]),
		})
		if err != nil {
			return nil, fmt.Errorf("cannot checkout tag: %s", err)
		}
	}
	if v, found := params["branch"]; found {
		err = worktree.Checkout(&git.CheckoutOptions{
			Branch: plumbing.NewRemoteReferenceName("origin", v[0]),
		})
		if err != nil {
			return nil, fmt.Errorf("cannot checkout branch: %s", err)
		}
	}

	paths, err := util.Glob(worktree.Filesystem, "*/*")
	if err != nil {
		return nil, fmt.Errorf("cannot list files: %s", err)
	}
	paths2, err := util.Glob(worktree.Filesystem, "*")
	if err != nil {
		return nil, fmt.Errorf("cannot list files: %s", err)
	}
	paths = append(paths, paths2...)


	fs = worktree.Filesystem
	files := map[string]string{}
	for _, path := range paths {
		info, err := fs.Stat(path)
		if err != nil {
			return nil, fmt.Errorf("cannot get file: %s", err)
		}

		if info.IsDir() {
			continue
		}

		f, err := fs.Open(path)
		if err != nil {
			return nil, fmt.Errorf("cannot get file: %s", err)
		}
		defer f.Close()

		content, err := ioutil.ReadAll(f)
		if err != nil {
			return nil, fmt.Errorf("cannot get file: %s", err)
		}

		files[path] = string(content)
	}

	return files, nil
}

func loadStackFromFS(root string) (map[string]string, error) {
	if !strings.HasSuffix(root, "/") {
		root = root + "/"
	}

	files := map[string]string{}

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		path = strings.TrimPrefix(path, root)
		if info.IsDir() {
			return nil
		}
		if strings.HasPrefix(path, "assets/") ||
			strings.HasPrefix(path, ".git/") {
			return nil
		}

		content, err := ioutil.ReadFile(filepath.Join(root, path))
		if err != nil {
			return fmt.Errorf("cannot get file: %s", err)
		}
		files[path] = string(content)

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("cannot walk %s: %s", root, err)
	}

	return files, nil
}
