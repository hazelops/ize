package config

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs/cloudwatchlogsiface"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecr/ecriface"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
)

type Project struct {
	TerraformVersion string `mapstructure:"terraform_version,omitempty"`
	AwsRegion        string `mapstructure:"aws_region,omitempty"`
	AwsProfile       string `mapstructure:"aws_profile,omitempty"`
	Namespace        string `mapstructure:"namespace,omitempty"`
	Env              string `mapstructure:"env,omitempty"`
	LogLevel         string `mapstructure:"log_level,omitempty"`
	PlainText        bool   `mapstructure:"plain_text_output,omitempty"`
	CustomPrompt     bool   `mapstructure:"custom_prompt,omitempty"`
	PreferRuntime    string `mapstructure:"prefer_runtime,omitempty"`
	Tag              string `mapstructure:"tag,omitempty"`
	DockerRegistry   string `mapstructure:"docker_registry,omitempty"`
	EndpointUrl      string `mapstructure:"endpoint_url,omitempty"`
	LocalStack       bool   `mapstructure:"localstack,omitempty"`
	SshPublicKey     string `mapstructure:"ssh_public_key,omitempty"`

	Home      string `mapstructure:"home,omitempty"`
	RootDir   string `mapstructure:"root_dir,omitempty"`
	InfraDir  string `mapstructure:"ize_dir,omitempty"`
	EnvDir    string `mapstructure:"env_dir,omitempty"`
	AppsPath  string `mapstructure:"apps_path,omitempty"`
	TFLog     string `mapstructure:"tf_log,omitempty"`
	TFLogPath string `mapstructure:"tf_log_path,omitempty"`

	Session   *session.Session
	AWSClient *awsClient

	Tunnel     *Tunnel                `mapstructure:",omitempty"`
	Terraform  map[string]*Terraform  `mapstructure:",omitempty"`
	Ecs        map[string]*Ecs        `mapstructure:",omitempty"`
	Serverless map[string]*Serverless `mapstructure:",omitempty"`
	Alias      map[string]*Alias      `mapstructure:",omitempty"`
}

type awsClient struct {
	S3Client             s3iface.S3API
	STSClient            stsiface.STSAPI
	IAMClient            iamiface.IAMAPI
	ECSClient            ecsiface.ECSAPI
	CloudWatchLogsClient cloudwatchlogsiface.CloudWatchLogsAPI
	SSMClient            ssmiface.SSMAPI
	ELBV2Client          elbv2iface.ELBV2API
	ECRClient            ecriface.ECRAPI
}

type Option func(*awsClient)

func WithS3Client(api s3iface.S3API) Option {
	return func(r *awsClient) {
		r.S3Client = api
	}
}

func WithSTSClient(api stsiface.STSAPI) Option {
	return func(r *awsClient) {
		r.STSClient = api
	}
}

func WithECSClient(api ecsiface.ECSAPI) Option {
	return func(r *awsClient) {
		r.ECSClient = api
	}
}

func WithSSMClient(api ssmiface.SSMAPI) Option {
	return func(r *awsClient) {
		r.SSMClient = api
	}
}

func WithELBV2Client(api elbv2iface.ELBV2API) Option {
	return func(r *awsClient) {
		r.ELBV2Client = api
	}
}

func WithCloudWatchLogsClient(api cloudwatchlogsiface.CloudWatchLogsAPI) Option {
	return func(r *awsClient) {
		r.CloudWatchLogsClient = api
	}
}

func WithIAMClient(api iamiface.IAMAPI) Option {
	return func(r *awsClient) {
		r.IAMClient = api
	}
}

func WithECRClient(api ecriface.ECRAPI) Option {
	return func(r *awsClient) {
		r.ECRClient = api
	}
}

func NewAWSClient(options ...Option) *awsClient {
	r := awsClient{}
	for _, opt := range options {
		opt(&r)
	}

	return &r
}

func (p *Project) SettingAWSClient(sess *session.Session) {
	p.AWSClient = NewAWSClient(
		WithS3Client(s3.New(sess)),
		WithSTSClient(sts.New(sess)),
		WithIAMClient(iam.New(sess)),
		WithECSClient(ecs.New(sess)),
		WithCloudWatchLogsClient(cloudwatchlogs.New(sess)),
		WithSSMClient(ssm.New(sess)),
		WithELBV2Client(elbv2.New(sess)),
		WithECRClient(ecr.New(sess)),
	)
}

func (p *Project) GetApps() map[string]*interface{} {
	apps := map[string]*interface{}{}

	for name, body := range p.Ecs {
		var v interface{}
		v = map[string]interface{}{
			"depends_on": body.DependsOn,
		}
		apps[name] = &v
	}

	for name, body := range p.Serverless {
		var v interface{}
		v = map[string]interface{}{
			"depends_on": body.DependsOn,
		}
		apps[name] = &v
	}

	for name, body := range p.Alias {
		var v interface{}
		v = map[string]interface{}{
			"depends_on": body.DependsOn,
		}
		apps[name] = &v
	}

	return apps
}

func (p *Project) GetStates() map[string]*interface{} {
	states := map[string]*interface{}{}

	for name, body := range p.Terraform {
		if name == "infra" {
			continue
		}
		var v interface{}
		v = map[string]interface{}{
			"depends_on": body.DependsOn,
		}
		states[name] = &v
	}

	return states
}
