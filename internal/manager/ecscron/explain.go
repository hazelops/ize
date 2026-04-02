package ecscron

import (
	"text/template"

	"github.com/hazelops/ize/internal/config"
)

func (m *Manager) Explain() error {
	m.prepare()
	return m.Project.Generate(ecsCronExplainTmpl, template.FuncMap{
		"app": func() config.EcsCron {
			return *m.App
		},
	})
}

var ecsCronExplainTmpl = `
# Create/Update EventBridge rule
aws events put-rule \
    --name {{.Env}}-{{app.Name}}-cron \
    --schedule-expression "{{app.Schedule}}" \
    --state ENABLED

# Set ECS task as target
aws events put-targets \
    --rule {{.Env}}-{{app.Name}}-cron \
    --targets '[{"Id":"{{app.Name}}","Arn":"arn:aws:ecs:{{.AwsRegion}}:<ACCOUNT_ID>:cluster/{{app.Cluster}}","EcsParameters":{"TaskDefinitionArn":"<TASK_DEF_ARN>","TaskCount":1,"LaunchType":"FARGATE"}}]'

# Run task manually
aws ecs run-task \
    --cluster {{app.Cluster}} \
    --task-definition {{.Env}}-{{app.Name}} \
    --launch-type FARGATE

# Remove targets and delete rule
aws events remove-targets --rule {{.Env}}-{{app.Name}}-cron --ids {{app.Name}}
aws events delete-rule --name {{.Env}}-{{app.Name}}-cron
`
