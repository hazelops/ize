package ecs

import (
	"text/template"

	"github.com/hazelops/ize/internal/config"
)

func (e *Manager) Explain() error {
	e.prepare()
	return e.Project.Generate(pushEcsAppTmpl, template.FuncMap{
		"app": func() config.Ecs {
			return *e.App
		},
	})
}

var pushEcsAppTmpl = `
# Authenticate the Docker CLI to registry
aws ecr get-login-password --region {{.AwsRegion}} | docker login --username AWS --password-stdin {{.DockerRegistry}}

# Create a repository
aws ecr describe-repositories --repository-names {{.Namespace}}-{{app.Name}} || \
aws ecr create-repository \
    --repository-name {{.Namespace}}-{{app.Name}} \
    --region {{.AwsRegion}}

# Push an image to Amazon ECR
docker push {{.DockerRegistry}}/{{.Namespace}}-{{app.Name}}:{{.Tag}} && \
docker push {{.DockerRegistry}}/{{.Namespace}}-{{app.Name}}:{{.Env}}-latest
`
