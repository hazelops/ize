package serverless

import (
	"text/template"

	"github.com/hazelops/ize/internal/config"
)

func (sls *Manager) Explain() error {
	sls.prepare()
	return sls.Project.Generate(upSLSAppTmpl, template.FuncMap{
		"app": func() config.Serverless {
			return *sls.App
		},
	})
}

var upSLSAppTmpl = `
# Change to the dir
cd {{app.Path}}

# Install the specified node version
nvm install {{app.NodeVersion}} 

# Switch to using the specified node version
nvm use {{app.NodeVersion}}

# Install dependencies
{{- if app.UseYarn}}
yarn install --save-dev

{{- if app.CreateDomain}}
# Create domain
yarn serverless create_domain \
	--region={{app.AwsRegion}} \
	--aws-profile={{app.AwsProfile}} \
	--stage={{.Project.Env}} \
	--verbose
{{- end}}

# Deploy serverless app
yarn serverless deploy \
	--config={{app.File}} \
	{{- if eq app.ServerlessVersion "3"}}
	--param="service={{app.Name}}" \
	{{- else}}
	--service={{app.Name}} \
	{{- end}}
	--region={{app.AwsRegion}} \
	--aws-profile={{app.AwsProfile}} \
	--stage={{.Project.Env}} \
	--verbose
{{- else}}
npm install --save-dev

{{- if app.CreateDomain}}
# Create domain
npx serverless create_domain \
	--region={{app.AwsRegion}} \
	--aws-profile={{app.AwsProfile}} \
	--stage={{.Project.Env}} \
	--verbose
{{- end}}

# Deploy serverless app
npx serverless deploy \
	--config={{app.File}} \
	{{- if eq app.ServerlessVersion "3"}}
	--param="service={{app.Name}}" \
	{{- else}}
	--service={{app.Name}} \
	{{- end}}
	--region={{app.AwsRegion}} \
	--aws-profile={{app.AwsProfile}} \
	--stage={{.Project.Env}} \
	--verbose
{{- end}}
`
