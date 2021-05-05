module github.com/gimlet-io/gimlet-stack

go 1.16

require (
	github.com/Masterminds/sprig/v3 v3.2.2 // indirect
	github.com/enescakir/emoji v1.0.0
	github.com/gimlet-io/gimlet-cli v0.8.0-rc1.0.20210428134552-eff0760ce3f3
	github.com/gimlet-io/gimletd v0.1.1-0.20210319125454-545e91c54f48 // indirect
	github.com/go-git/go-billy/v5 v5.0.0 // indirect
	github.com/go-git/go-git/v5 v5.3.0 // indirect
	github.com/stretchr/testify v1.7.0 // indirect
	github.com/urfave/cli/v2 v2.3.0
	github.com/whilp/git-urls v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect
)

replace (
	// https://github.com/helm/helm/issues/9354
	github.com/docker/distribution => github.com/docker/distribution v0.0.0-20191216044856-a8371794149d
	github.com/docker/docker => github.com/moby/moby v17.12.0-ce-rc1.0.20200618181300-9dc6525e6118+incompatible

	github.com/go-git/go-git/v5 => github.com/gimlet-io/go-git/v5 v5.2.1-0.20210122134038-45142aa695dd
)
