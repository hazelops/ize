package commands

import (
	_ "embed"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/golang/mock/gomock"
	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/pkg/mocks"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

//go:generate mockgen -package=mocks -destination ../../pkg/mocks/mock_s3.go github.com/aws/aws-sdk-go/service/s3/s3iface S3API
//go:generate mockgen -package=mocks -destination ../../pkg/mocks/mock_sts.go github.com/aws/aws-sdk-go/service/sts/stsiface STSAPI

//go:embed testdata/tfenv_valid.toml
var tfenvToml string

func TestTfenv(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		wantErr        bool
		wantBackend    string
		wantTfvars     string
		withConfigFile bool
		env            map[string]string
		mockS3Client   func(m *mocks.MockS3API)
		mockSTSClient  func(m *mocks.MockSTSAPI)
		withECS        bool
	}{
		{
			name:           "success (only config file)",
			args:           []string{"gen", "tfenv"},
			env:            map[string]string{"ENV": "test", "AWS_PROFILE": "test"},
			withConfigFile: true,
			wantErr:        false,
			withECS:        true,
			mockS3Client: func(m *mocks.MockS3API) {
				m.EXPECT().HeadBucket(gomock.Any()).Return(nil, nil).AnyTimes()
				m.EXPECT().HeadObject(gomock.Any()).Return(nil, nil).AnyTimes()
			},
			mockSTSClient: func(m *mocks.MockSTSAPI) {},
			wantBackend: `provider "aws" {
  profile = var.aws_profile
  region  = var.aws_region
  default_tags {
    tags = {
      env       = "test"
      namespace = "testnut"
      terraform = "true"
    }
  }
}

terraform {
  backend "s3" {
    bucket         = "testnut-tf-state"
    key            = "test/terraform.tfstate"
    region         = "us-east-1"
    profile        = "test"
    dynamodb_table = "tf-state-lock"
  }
}
`,
			wantTfvars: `env               = "test"
aws_profile       = "test"
aws_region        = "us-east-1"
ec2_key_pair_name = "test-testnut"
docker_image_tag  = "test"
ssh_public_key    =
docker_registry   = "0.dkr.ecr.us-east-1.amazonaws.com"
namespace         = "testnut"
root_domain_name  = "examples.ize.sh"
`,
		},
		{
			name:           "success (flags and config file)",
			args:           []string{"gen", "tfenv", "--terraform-state-bucket-name=test"},
			env:            map[string]string{"ENV": "test", "AWS_PROFILE": "test"},
			withConfigFile: true,
			wantErr:        false,
			mockS3Client: func(m *mocks.MockS3API) {
				m.EXPECT().HeadBucket(gomock.Any()).Return(nil, nil).AnyTimes()
				m.EXPECT().HeadObject(gomock.Any()).Return(nil, nil).AnyTimes()
			},
			mockSTSClient: func(m *mocks.MockSTSAPI) {},
			withECS:       true,
			wantBackend: `provider "aws" {
  profile = var.aws_profile
  region  = var.aws_region
  default_tags {
    tags = {
      env       = "test"
      namespace = "testnut"
      terraform = "true"
    }
  }
}

terraform {
  backend "s3" {
    bucket         = "test"
    key            = "test/terraform.tfstate"
    region         = "us-east-1"
    profile        = "test"
    dynamodb_table = "tf-state-lock"
  }
}
`,
			wantTfvars: `env               = "test"
aws_profile       = "test"
aws_region        = "us-east-1"
ec2_key_pair_name = "test-testnut"
docker_image_tag  = "test"
ssh_public_key    =
docker_registry   = "0.dkr.ecr.us-east-1.amazonaws.com"
namespace         = "testnut"
root_domain_name  = "examples.ize.sh"
`,
		},
		{
			name:           "success (flags, env and config file)",
			args:           []string{"gen", "tfenv", "--terraform-state-bucket-name=test"},
			env:            map[string]string{"ENV": "test", "AWS_PROFILE": "test", "IZE_TERRAFORM__INFRA__ROOT_DOMAIN_NAME": "test"},
			withConfigFile: true,
			wantErr:        false,
			mockS3Client: func(m *mocks.MockS3API) {
				m.EXPECT().HeadBucket(gomock.Any()).Return(nil, nil).AnyTimes()
				m.EXPECT().HeadObject(gomock.Any()).Return(nil, nil).AnyTimes()
			},
			withECS:       true,
			mockSTSClient: func(m *mocks.MockSTSAPI) {},
			wantBackend: `provider "aws" {
  profile = var.aws_profile
  region  = var.aws_region
  default_tags {
    tags = {
      env       = "test"
      namespace = "testnut"
      terraform = "true"
    }
  }
}

terraform {
  backend "s3" {
    bucket         = "test"
    key            = "test/terraform.tfstate"
    region         = "us-east-1"
    profile        = "test"
    dynamodb_table = "tf-state-lock"
  }
}
`,
			wantTfvars: `env               = "test"
aws_profile       = "test"
aws_region        = "us-east-1"
ec2_key_pair_name = "test-testnut"
docker_image_tag  = "test"
ssh_public_key    =
docker_registry   = "0.dkr.ecr.us-east-1.amazonaws.com"
namespace         = "testnut"
root_domain_name  = "test"
`,
		},
		{
			name:    "success (flags and env)",
			args:    []string{"--aws-region", "us-east-1", "--namespace", "testnut", "gen", "tfenv"},
			env:     map[string]string{"ENV": "test", "AWS_PROFILE": "test", "IZE_TERRAFORM__INFRA__ROOT_DOMAIN_NAME": "test"},
			wantErr: false,
			mockS3Client: func(m *mocks.MockS3API) {
				m.EXPECT().HeadBucket(gomock.Any()).Return(nil, awserr.New(s3.ErrCodeNoSuchKey, "message", nil)).Times(1)
				m.EXPECT().HeadObject(gomock.Any()).Return(nil, nil).AnyTimes()
			},
			mockSTSClient: func(m *mocks.MockSTSAPI) {
				m.EXPECT().GetCallerIdentity(gomock.Any()).Return(&sts.GetCallerIdentityOutput{
					Account: aws.String("0"),
				}, nil).Times(1)
			},
			wantBackend: `provider "aws" {
  profile = var.aws_profile
  region  = var.aws_region
  default_tags {
    tags = {
      env       = "test"
      namespace = "testnut"
      terraform = "true"
    }
  }
}

terraform {
  backend "s3" {
    bucket         = "testnut-0-tf-state"
    key            = "test/terraform.tfstate"
    region         = "us-east-1"
    profile        = "test"
    dynamodb_table = "tf-state-lock"
  }
}
`,
			wantTfvars: `env               = "test"
aws_profile       = "test"
aws_region        = "us-east-1"
ec2_key_pair_name = "test-testnut"
ssh_public_key    =
namespace         = "testnut"
`,
		},
		{
			name:    "success (only flags)",
			args:    []string{"-e=test", "-r=us-east-1", "-p=test", "-n=testnut", "gen", "tfenv", "--terraform-state-bucket-name=test"},
			wantErr: false,
			mockS3Client: func(m *mocks.MockS3API) {
				m.EXPECT().HeadObject(gomock.Any()).Return(nil, nil).AnyTimes()
			},
			mockSTSClient: func(m *mocks.MockSTSAPI) {},
			wantBackend: `provider "aws" {
  profile = var.aws_profile
  region  = var.aws_region
  default_tags {
    tags = {
      env       = "test"
      namespace = "testnut"
      terraform = "true"
    }
  }
}

terraform {
  backend "s3" {
    bucket         = "test"
    key            = "test/terraform.tfstate"
    region         = "us-east-1"
    profile        = "test"
    dynamodb_table = "tf-state-lock"
  }
}
`,
			wantTfvars: `env               = "test"
aws_profile       = "test"
aws_region        = "us-east-1"
ec2_key_pair_name = "test-testnut"
ssh_public_key    =
namespace         = "testnut"
`,
		},
		{
			name:    "success (only env)",
			args:    []string{"gen", "tfenv"},
			env:     map[string]string{"ENV": "test", "AWS_PROFILE": "test", "NAMESPACE": "testnut", "AWS_REGION": "us-west-2", "IZE_TERRAFORM__INFRA__ROOT_DOMAIN_NAME": "test"},
			wantErr: false,
			mockS3Client: func(m *mocks.MockS3API) {
				m.EXPECT().HeadBucket(gomock.Any()).Return(nil, awserr.New(s3.ErrCodeNoSuchBucket, "message", nil)).Times(1)
				m.EXPECT().HeadObject(gomock.Any()).Return(nil, nil).AnyTimes()
			},
			mockSTSClient: func(m *mocks.MockSTSAPI) {
				m.EXPECT().GetCallerIdentity(gomock.Any()).Return(&sts.GetCallerIdentityOutput{
					Account: aws.String("0"),
				}, nil).Times(1)
			},
			wantBackend: `provider "aws" {
  profile = var.aws_profile
  region  = var.aws_region
  default_tags {
    tags = {
      env       = "test"
      namespace = "testnut"
      terraform = "true"
    }
  }
}

terraform {
  backend "s3" {
    bucket         = "testnut-0-tf-state"
    key            = "test/terraform.tfstate"
    region         = "us-west-2"
    profile        = "test"
    dynamodb_table = "tf-state-lock"
  }
}
`,
			wantTfvars: `env               = "test"
aws_profile       = "test"
aws_region        = "us-west-2"
ec2_key_pair_name = "test-testnut"
ssh_public_key    =
namespace         = "testnut"
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer resetEnv(os.Environ())
			viper.Reset()
			os.Unsetenv("IZE_CONFIG_FILE")
			// Set env
			for k, v := range tt.env {
				t.Setenv(k, v)
			}

			temp, err := os.MkdirTemp("", "test")
			if err != nil {
				t.Error(err)
				return
			}
			err = os.Chdir(temp)
			if err != nil {
				t.Error(err)
				return
			}
			err = os.MkdirAll(filepath.Join(temp, ".ize", "env", "test"), 0777)
			if err != nil {
				t.Error(err)
				return
			}

			if tt.withConfigFile {
				setConfigFile(filepath.Join(temp, "ize.toml"), tfenvToml, t)
				if tt.withECS {
					f, err := os.OpenFile(filepath.Join(temp, "ize.toml"), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
					if err != nil {
						panic(err)
					}
					defer f.Close()
					if _, err = f.WriteString("\n\n[ecs.squibby]\ntimeout = 0"); err != nil {
						panic(err)
					}
				}
			}

			t.Setenv("HOME", temp)

			err = os.MkdirAll(filepath.Join(temp, ".ssh"), 0777)
			if err != nil {
				t.Error(err)
				return
			}
			_, err = makeSSHKeyPair(filepath.Join(temp, ".ssh", "id_rsa.pub"), filepath.Join(temp, ".ssh", "id_rsa"))
			if err != nil {
				t.Error(err)
				return
			}

			key, err := os.ReadFile(filepath.Join(temp, ".ssh", "id_rsa.pub"))
			if err != nil {
				t.Error(err)
				return
			}

			tt.wantTfvars = strings.ReplaceAll(tt.wantTfvars, "ssh_public_key    =", fmt.Sprintf("ssh_public_key    = \"%s\"", strings.TrimSpace(string(key))))

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockS3API := mocks.NewMockS3API(ctrl)
			tt.mockS3Client(mockS3API)

			mockSTSAPI := mocks.NewMockSTSAPI(ctrl)
			tt.mockSTSClient(mockSTSAPI)

			cfg := new(config.Project)
			cmd := newRootCmd(cfg)

			cmd.SetArgs(tt.args)
			cmd.PersistentFlags().ParseErrorsWhitelist.UnknownFlags = true
			err = cmd.PersistentFlags().Parse(tt.args)
			if err != nil {
				t.Error(err)
				return
			}

			cmd.PersistentFlags().VisitAll(func(f *pflag.Flag) {
				if len(f.Value.String()) != 0 {
					_ = viper.BindPFlag(strings.ReplaceAll(f.Name, "-", "_"), cmd.PersistentFlags().Lookup(f.Name))
				}
			})

			config.InitConfig()

			err = cfg.GetTestConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("get config error = %v, wantErr %v", err, tt.wantErr)
				os.Exit(1)
			}
			cfg.AWSClient = config.NewAWSClient(
				config.WithS3Client(mockS3API),
				config.WithSTSClient(mockSTSAPI),
			)

			err = cmd.Execute()
			if (err != nil) != tt.wantErr {
				t.Errorf("ize gen tfenv error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			file, err := os.ReadFile(filepath.Join(temp, ".ize", "env", "test", "backend.tf"))
			if err != nil {
				t.Error(err)
				return
			}

			if !reflect.DeepEqual(string(file), tt.wantBackend) {
				t.Errorf("backend.tf = %v, want %v", string(file), tt.wantBackend)
			}

			file, err = os.ReadFile(filepath.Join(temp, ".ize", "env", "test", "terraform.tfvars"))
			if err != nil {
				t.Error(err)
				return
			}

			if !reflect.DeepEqual(string(file), tt.wantTfvars) {
				t.Errorf("terraform.tfvars = %v, want %v", string(file), tt.wantTfvars)
				return
			}

			// Unset env
			for k, _ := range tt.env {
				os.Unsetenv(k)
			}
		})
	}
}
