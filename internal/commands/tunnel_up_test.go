package commands

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"github.com/hazelops/ize/internal/config"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestUpOptions_getSSHCommandArgs(t *testing.T) {
	type fields struct {
		Config                *config.Project
		PrivateKeyFile        string
		PublicKeyFile         string
		BastionHostID         string
		ForwardHost           []string
		StrictHostKeyChecking bool
	}
	type args struct {
		sshConfigPath string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		{name: "success",
			fields: fields{
				Config:                &config.Project{},
				PrivateKeyFile:        "",
				PublicKeyFile:         "",
				BastionHostID:         "i-xxxxxxxxxxx",
				ForwardHost:           nil,
				StrictHostKeyChecking: false,
			},
			args: args{sshConfigPath: "./test/ssh.config"},
			want: []string{"-M", "-t", "-S", "bastion.sock", "-fN", "ubuntu@i-xxxxxxxxxxx", "-F", "./test/ssh.config"},
		},
		{name: "success with StrictHostKeyChecking",
			fields: fields{
				Config:                &config.Project{},
				PrivateKeyFile:        "",
				PublicKeyFile:         "",
				BastionHostID:         "i-xxxxxxxxxxx",
				ForwardHost:           nil,
				StrictHostKeyChecking: true,
			},
			args: args{sshConfigPath: "./test/ssh.config"},
			want: []string{"-M", "-t", "-S", "bastion.sock", "-fN", "-o", "StrictHostKeyChecking=no", "ubuntu@i-xxxxxxxxxxx", "-F", "./test/ssh.config"},
		},
		{name: "success with private key file",
			fields: fields{
				Config:                &config.Project{},
				PrivateKeyFile:        fmt.Sprintf("%s/.ssh/id_rsa", homeDir()),
				PublicKeyFile:         "",
				BastionHostID:         "i-XXXXXXXXXXXXXXXXX",
				ForwardHost:           nil,
				StrictHostKeyChecking: false,
			},
			args: args{sshConfigPath: "./test/ssh.config"},
			want: []string{"-M", "-t", "-S", "bastion.sock", "-fN", "ubuntu@i-XXXXXXXXXXXXXXXXX", "-F", "./test/ssh.config", "-i", fmt.Sprintf("%s/.ssh/id_rsa", homeDir())},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &TunnelUpOptions{
				Config:                tt.fields.Config,
				PrivateKeyFile:        tt.fields.PrivateKeyFile,
				PublicKeyFile:         tt.fields.PublicKeyFile,
				BastionHostID:         tt.fields.BastionHostID,
				ForwardHost:           tt.fields.ForwardHost,
				StrictHostKeyChecking: tt.fields.StrictHostKeyChecking,
			}
			if got := o.getSSHCommandArgs(tt.args.sshConfigPath); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getSSHCommandArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func homeDir() string {
	dir, _ := os.UserHomeDir()
	return dir
}

type mockSSM struct {
	ssmiface.SSMAPI
	response string
	err      error
}

func (sp mockSSM) GetParameter(*ssm.GetParameterInput) (*ssm.GetParameterOutput, error) {
	output := new(ssm.GetParameterOutput)
	output.Parameter = &ssm.Parameter{Value: aws.String(sp.response)}
	return output, sp.err
}

func Test_getTerraformOutput(t *testing.T) {
	type args struct {
		wr  *SSMWrapper
		env string
	}
	tests := []struct {
		name    string
		args    args
		want    terraformOutput
		wantErr bool
	}{
		{
			name: "success", args: args{
				wr: &SSMWrapper{Api: mockSSM{
					SSMAPI:   nil,
					response: "ew0KICAiYmFzdGlvbl9pbnN0YW5jZV9pZCI6IHsNCiAgICAic2Vuc2l0aXZlIjogZmFsc2UsDQogICAgInR5cGUiOiAic3RyaW5nIiwNCiAgICAidmFsdWUiOiAiaS1YWFhYWFhYWFhYWFhYWFhYWCINCiAgfSwNCiAgImNtZCI6IHsNCiAgICAic2Vuc2l0aXZlIjogZmFsc2UsDQogICAgInR5cGUiOiBbDQogICAgICAib2JqZWN0IiwNCiAgICAgIHsNCiAgICAgICAgInR1bm5lbCI6IFsNCiAgICAgICAgICAib2JqZWN0IiwNCiAgICAgICAgICB7DQogICAgICAgICAgICAiZG93biI6ICJzdHJpbmciLA0KICAgICAgICAgICAgInN0YXR1cyI6ICJzdHJpbmciLA0KICAgICAgICAgICAgInVwIjogInN0cmluZyINCiAgICAgICAgICB9DQogICAgICAgIF0NCiAgICAgIH0NCiAgICBdLA0KICAgICJ2YWx1ZSI6IHsNCiAgICAgICJ0dW5uZWwiOiB7DQogICAgICAgICJkb3duIjogInNzaCAtUyBiYXN0aW9uLnNvY2sgLU8gZXhpdCB1YnVudHVAaS1YWFhYWFhYWFhYWFhYWFhYWCAiLA0KICAgICAgICAic3RhdHVzIjogInNzaCAtUyBiYXN0aW9uLnNvY2sgLU8gY2hlY2sgdWJ1bnR1QGktWFhYWFhYWFhYWFhYWFhYWFgiLA0KICAgICAgICAidXAiOiAic3NoIC1NIC1TIGJhc3Rpb24uc29jayAtZk5UIHVidW50dUBpLVhYWFhYWFhYWFhYWFhYWFhYICINCiAgICAgIH0NCiAgICB9DQogIH0sDQogICJzc2hfZm9yd2FyZF9jb25maWciOiB7DQogICAgInNlbnNpdGl2ZSI6IGZhbHNlLA0KICAgICJ0eXBlIjogWw0KICAgICAgInR1cGxlIiwNCiAgICAgIFsNCiAgICAgICAgInN0cmluZyIsDQogICAgICAgICJzdHJpbmciLA0KICAgICAgICAic3RyaW5nIiwNCiAgICAgICAgInN0cmluZyIsDQogICAgICAgICJzdHJpbmciDQogICAgICBdDQogICAgXSwNCiAgICAidmFsdWUiOiBbDQogICAgICAiIyBTU0ggb3ZlciBTZXNzaW9uIE1hbmFnZXIiLA0KICAgICAgImhvc3QgaS0qIG1pLSoiLA0KICAgICAgIlNlcnZlckFsaXZlSW50ZXJ2YWwgMTgwIiwNCiAgICAgICJQcm94eUNvbW1hbmQgc2ggLWMgXCJhd3Mgc3NtIHN0YXJ0LXNlc3Npb24gLS10YXJnZXQgJWggLS1kb2N1bWVudC1uYW1lIEFXUy1TdGFydFNTSFNlc3Npb24gLS1wYXJhbWV0ZXJzICdwb3J0TnVtYmVyPSVwJ1wiXG4iLA0KICAgICAgIkxvY2FsRm9yd2FyZCAzMjA4NCB0ZXN0LnRlc3QudGVzdDo4MCINCiAgICBdDQogIH0sDQogICJ2cGNfcHJpdmF0ZV9zdWJuZXRzIjogew0KICAgICJzZW5zaXRpdmUiOiBmYWxzZSwNCiAgICAidHlwZSI6IFsNCiAgICAgICJ0dXBsZSIsDQogICAgICBbDQogICAgICAgICJzdHJpbmciDQogICAgICBdDQogICAgXSwNCiAgICAidmFsdWUiOiBbDQogICAgICAic3VibmV0LVhYWFhYWFhYWFhYWFhYWFhYIg0KICAgIF0NCiAgfSwNCiAgInZwY19wdWJsaWNfc3VibmV0cyI6IHsNCiAgICAic2Vuc2l0aXZlIjogZmFsc2UsDQogICAgInR5cGUiOiBbDQogICAgICAidHVwbGUiLA0KICAgICAgWw0KICAgICAgICAic3RyaW5nIg0KICAgICAgXQ0KICAgIF0sDQogICAgInZhbHVlIjogWw0KICAgICAgInN1Ym5ldC1YWFhYWFhYWFhYWFhYWFhYWCINCiAgICBdDQogIH0NCn0=",
					err:      nil,
				}},
				env: "testnut",
			},
			want: terraformOutput{
				BastionInstanceID: struct {
					Value string `json:"value,omitempty"`
				}{Value: "i-XXXXXXXXXXXXXXXXX"},
				SSHForwardConfig: struct {
					Value []string `json:"value,omitempty"`
				}{Value: []string{"# SSH over Session Manager", "host i-* mi-*", "ServerAliveInterval 180", "ProxyCommand sh -c \"aws ssm start-session --target %h --document-name AWS-StartSSHSession --parameters 'portNumber=%p'\"\n", "LocalForward 32084 test.test.test:80"}},
			},
			wantErr: false,
		},
		{
			name: "incorrect json", args: args{
				wr: &SSMWrapper{Api: mockSSM{
					SSMAPI:   nil,
					response: "ICAiYmFzdGlvbl9pbnN0YW5jZV9pZCI6IHsNCiAgICAic2Vuc2l0aXZlIjogZmFsc2UsDQogICAgInR5cGUiOiAic3RyaW5nIiwNCiAgICAidmFsdWUiOiAiaS1YWFhYWFhYWFhYWFhYWFhYWCINCiAgfSwNCiAgImNtZCI6IHsNCiAgICAic2Vuc2l0aXZlIjogZmFsc2UsDQogICAgInR5cGUiOiBbDQogICAgICAib2JqZWN0IiwNCiAgICAgIHsNCiAgICAgICAgInR1bm5lbCI6IFsNCiAgICAgICAgICAib2JqZWN0IiwNCiAgICAgICAgICB7DQogICAgICAgICAgICAiZG93biI6ICJzdHJpbmciLA0KICAgICAgICAgICAgInN0YXR1cyI6ICJzdHJpbmciLA0KICAgICAgICAgICAgInVwIjogInN0cmluZyINCiAgICAgICAgICB9DQogICAgICAgIF0NCiAgICAgIH0NCiAgICBdLA0KICAgICJ2YWx1ZSI6IHsNCiAgICAgICJ0dW5uZWwiOiB7DQogICAgICAgICJkb3duIjogInNzaCAtUyBiYXN0aW9uLnNvY2sgLU8gZXhpdCB1YnVudHVAaS1YWFhYWFhYWFhYWFhYWFhYWCAiLA0KICAgICAgICAic3RhdHVzIjogInNzaCAtUyBiYXN0aW9uLnNvY2sgLU8gY2hlY2sgdWJ1bnR1QGktWFhYWFhYWFhYWFhYWFhYWFgiLA0KICAgICAgICAidXAiOiAic3NoIC1NIC1TIGJhc3Rpb24uc29jayAtZk5UIHVidW50dUBpLVhYWFhYWFhYWFhYWFhYWFhYICINCiAgICAgIH0NCiAgICB9DQogIH0sDQogICJzc2hfZm9yd2FyZF9jb25maWciOiB7DQogICAgInNlbnNpdGl2ZSI6IGZhbHNlLA0KICAgICJ0eXBlIjogWw0KICAgICAgInR1cGxlIiwNCiAgICAgIFsNCiAgICAgICAgInN0cmluZyIsDQogICAgICAgICJzdHJpbmciLA0KICAgICAgICAic3RyaW5nIiwNCiAgICAgICAgInN0cmluZyIsDQogICAgICAgICJzdHJpbmciDQogICAgICBdDQogICAgXSwNCiAgICAidmFsdWUiOiBbDQogICAgICAiIyBTU0ggb3ZlciBTZXNzaW9uIE1hbmFnZXIiLA0KICAgICAgImhvc3QgaS0qIG1pLSoiLA0KICAgICAgIlNlcnZlckFsaXZlSW50ZXJ2YWwgMTgwIiwNCiAgICAgICJQcm94eUNvbW1hbmQgc2ggLWMgXCJhd3Mgc3NtIHN0YXJ0LXNlc3Npb24gLS10YXJnZXQgJWggLS1kb2N1bWVudC1uYW1lIEFXUy1TdGFydFNTSFNlc3Npb24gLS1wYXJhbWV0ZXJzICdwb3J0TnVtYmVyPSVwJ1wiXG4iLA0KICAgICAgIkxvY2FsRm9yd2FyZCAzMjA4NCB0ZXN0LnRlc3QudGVzdDo4MCINCiAgICBdDQogIH0sDQogICJ2cGNfcHJpdmF0ZV9zdWJuZXRzIjogew0KICAgICJzZW5zaXRpdmUiOiBmYWxzZSwNCiAgICAidHlwZSI6IFsNCiAgICAgICJ0dXBsZSIsDQogICAgICBbDQogICAgICAgICJzdHJpbmciDQogICAgICBdDQogICAgXSwNCiAgICAidmFsdWUiOiBbDQogICAgICAic3VibmV0LVhYWFhYWFhYWFhYWFhYWFhYIg0KICAgIF0NCiAgfSwNCiAgInZwY19wdWJsaWNfc3VibmV0cyI6IHsNCiAgICAic2Vuc2l0aXZlIjogZmFsc2UsDQogICAgInR5cGUiOiBbDQogICAgICAidHVwbGUiLA0KICAgICAgWw0KICAgICAgICAic3RyaW5nIg0KICAgICAgXQ0KICAgIF0sDQogICAgInZhbHVlIjogWw0KICAgICAgInN1Ym5ldC1YWFhYWFhYWFhYWFhYWFhYWCINCiAgICBdDQogIH0NCn0",
					err:      nil,
				}},
				env: "testnut",
			},
			want:    terraformOutput{},
			wantErr: true,
		},
		{
			name: "incorrect base64", args: args{
				wr: &SSMWrapper{Api: mockSSM{
					SSMAPI:   nil,
					response: "0a2a-67c5-7a30-a153-b88f-81ac-8898-b766",
					err:      nil,
				}},
				env: "testnut",
			},
			want:    terraformOutput{},
			wantErr: true,
		},
		{
			name: "incorrect path", args: args{
				wr: &SSMWrapper{Api: mockSSM{
					SSMAPI:   nil,
					response: "",
					err:      awserr.New(ssm.ErrCodeParameterNotFound, "", nil),
				}},
				env: "incorrectPath",
			},
			want:    terraformOutput{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getTerraformOutput(tt.args.wr, tt.args.env)
			if (err != nil) != tt.wantErr {
				t.Errorf("getTerraformOutput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getTerraformOutput() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_writeSSHConfigFromConfig(t *testing.T) {
	tmp := os.TempDir()

	type args struct {
		forwardHost []string
		dir         string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				forwardHost: []string{"11.11.11.11:22:36002", "11.11.11.11:23:36001"},
				dir:         tmp,
			},
			want:    "# SSH over Session Manager\nhost i-* mi-*\nServerAliveInterval 180\nProxyCommand sh -c \"aws ssm start-session --target %h --document-name AWS-StartSSHSession --parameters 'portNumber=%p'\"\n\nLocalForward 36002 11.11.11.11:22\nLocalForward 36001 11.11.11.11:23\n\n",
			wantErr: false,
		},
		{
			name: "incorrect forward host 1",
			args: args{
				forwardHost: []string{"11.11.11.11", "11.11.11.11:23:36001"},
				dir:         tmp,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "incorrect forward host 2",
			args: args{
				forwardHost: []string{"11.11.11.11:22:", "11.11.11.11:23:36001"},
				dir:         tmp,
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := writeSSHConfigFromConfig(tt.args.forwardHost, tt.args.dir); (err != nil) != tt.wantErr {
				t.Errorf("writeSSHConfigFromConfig() error = %v, wantErr %v", err, tt.wantErr)
				t.SkipNow()
			}

			file, err := os.ReadFile(filepath.Join(tmp, "ssh.config"))
			if err != nil {
				t.Error("writeSSHConfigFromConfig(): ssh.config not found")
			}

			fmt.Println(string(file))
			if !reflect.DeepEqual(string(file), tt.want) {
				t.Errorf("writeSSHConfigFromConfig() = %v, want %v", string(file), tt.want)
			}
		})
	}
}

func Test_writeSSHConfigFromSSM(t *testing.T) {
	tmp := os.TempDir()

	type args struct {
		wr  *SSMWrapper
		env string
		dir string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   []string
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				wr: &SSMWrapper{Api: mockSSM{
					SSMAPI:   nil,
					response: "ew0KICAiYmFzdGlvbl9pbnN0YW5jZV9pZCI6IHsNCiAgICAic2Vuc2l0aXZlIjogZmFsc2UsDQogICAgInR5cGUiOiAic3RyaW5nIiwNCiAgICAidmFsdWUiOiAiaS1YWFhYWFhYWFhYWFhYWFhYWCINCiAgfSwNCiAgImNtZCI6IHsNCiAgICAic2Vuc2l0aXZlIjogZmFsc2UsDQogICAgInR5cGUiOiBbDQogICAgICAib2JqZWN0IiwNCiAgICAgIHsNCiAgICAgICAgInR1bm5lbCI6IFsNCiAgICAgICAgICAib2JqZWN0IiwNCiAgICAgICAgICB7DQogICAgICAgICAgICAiZG93biI6ICJzdHJpbmciLA0KICAgICAgICAgICAgInN0YXR1cyI6ICJzdHJpbmciLA0KICAgICAgICAgICAgInVwIjogInN0cmluZyINCiAgICAgICAgICB9DQogICAgICAgIF0NCiAgICAgIH0NCiAgICBdLA0KICAgICJ2YWx1ZSI6IHsNCiAgICAgICJ0dW5uZWwiOiB7DQogICAgICAgICJkb3duIjogInNzaCAtUyBiYXN0aW9uLnNvY2sgLU8gZXhpdCB1YnVudHVAaS1YWFhYWFhYWFhYWFhYWFhYWCAiLA0KICAgICAgICAic3RhdHVzIjogInNzaCAtUyBiYXN0aW9uLnNvY2sgLU8gY2hlY2sgdWJ1bnR1QGktWFhYWFhYWFhYWFhYWFhYWFgiLA0KICAgICAgICAidXAiOiAic3NoIC1NIC1TIGJhc3Rpb24uc29jayAtZk5UIHVidW50dUBpLVhYWFhYWFhYWFhYWFhYWFhYICINCiAgICAgIH0NCiAgICB9DQogIH0sDQogICJzc2hfZm9yd2FyZF9jb25maWciOiB7DQogICAgInNlbnNpdGl2ZSI6IGZhbHNlLA0KICAgICJ0eXBlIjogWw0KICAgICAgInR1cGxlIiwNCiAgICAgIFsNCiAgICAgICAgInN0cmluZyIsDQogICAgICAgICJzdHJpbmciLA0KICAgICAgICAic3RyaW5nIiwNCiAgICAgICAgInN0cmluZyIsDQogICAgICAgICJzdHJpbmciDQogICAgICBdDQogICAgXSwNCiAgICAidmFsdWUiOiBbDQogICAgICAiIyBTU0ggb3ZlciBTZXNzaW9uIE1hbmFnZXIiLA0KICAgICAgImhvc3QgaS0qIG1pLSoiLA0KICAgICAgIlNlcnZlckFsaXZlSW50ZXJ2YWwgMTgwIiwNCiAgICAgICJQcm94eUNvbW1hbmQgc2ggLWMgXCJhd3Mgc3NtIHN0YXJ0LXNlc3Npb24gLS10YXJnZXQgJWggLS1kb2N1bWVudC1uYW1lIEFXUy1TdGFydFNTSFNlc3Npb24gLS1wYXJhbWV0ZXJzICdwb3J0TnVtYmVyPSVwJ1wiXG4iLA0KICAgICAgIkxvY2FsRm9yd2FyZCAzMjA4NCB0ZXN0LnRlc3QudGVzdDo4MCINCiAgICBdDQogIH0sDQogICJ2cGNfcHJpdmF0ZV9zdWJuZXRzIjogew0KICAgICJzZW5zaXRpdmUiOiBmYWxzZSwNCiAgICAidHlwZSI6IFsNCiAgICAgICJ0dXBsZSIsDQogICAgICBbDQogICAgICAgICJzdHJpbmciDQogICAgICBdDQogICAgXSwNCiAgICAidmFsdWUiOiBbDQogICAgICAic3VibmV0LVhYWFhYWFhYWFhYWFhYWFhYIg0KICAgIF0NCiAgfSwNCiAgInZwY19wdWJsaWNfc3VibmV0cyI6IHsNCiAgICAic2Vuc2l0aXZlIjogZmFsc2UsDQogICAgInR5cGUiOiBbDQogICAgICAidHVwbGUiLA0KICAgICAgWw0KICAgICAgICAic3RyaW5nIg0KICAgICAgXQ0KICAgIF0sDQogICAgInZhbHVlIjogWw0KICAgICAgInN1Ym5ldC1YWFhYWFhYWFhYWFhYWFhYWCINCiAgICBdDQogIH0NCn0=",
					err:      nil,
				}},
				env: "test",
				dir: tmp,
			},
			want:    "i-XXXXXXXXXXXXXXXXX",
			want1:   []string{"test.test.test:80:32084"},
			wantErr: false,
		},
		{
			name: "SSM error",
			args: args{
				wr: &SSMWrapper{Api: mockSSM{
					SSMAPI: nil,
					err:    awserr.New(ssm.ErrCodeParameterNotFound, "", nil),
				}},
				env: "test",
				dir: tmp,
			},
			want:    "",
			want1:   []string{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := writeSSHConfigFromSSM(tt.args.wr, tt.args.env, tt.args.dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("writeSSHConfigFromSSM() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("writeSSHConfigFromSSM() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("writeSSHConfigFromSSM() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_setAWSCredentials(t *testing.T) {
	type args struct {
		sess *session.Session
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "error", args: args{sess: getSession(true)}, wantErr: true},
		{name: "success", args: args{sess: getSession(false)}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := setAWSCredentials(tt.args.sess); (err != nil) != tt.wantErr {
				t.Errorf("setAWSCredentials() error = %v, wantErr %v", err, tt.wantErr)
			}
			if os.Getenv("AWS_SECRET_ACCESS_KEY") == "test" && os.Getenv("AWS_ACCESS_KEY_ID") == "test" {
				return
			}
		})
	}
}

func getSession(incorrect bool) *session.Session {
	if incorrect {
		return session.Must(session.NewSession(aws.NewConfig().WithCredentials(credentials.NewSharedCredentials(filepath.Join("incorrect", "credentials"), "test"))))
	}

	tmp := os.TempDir()

	testCreds := `[test]
aws_access_key_id = test
aws_secret_access_key = test
`

	f, _ := os.Create(filepath.Join(tmp, "credentials"))
	_, _ = io.WriteString(f, testCreds)

	return session.Must(session.NewSession(aws.NewConfig().WithCredentials(credentials.NewSharedCredentials(filepath.Join(tmp, "credentials"), "test"))))
}

func TestUpOptions_runSSH(t *testing.T) {
	type fields struct {
		Config                *config.Project
		PrivateKeyFile        string
		PublicKeyFile         string
		BastionHostID         string
		ForwardHost           []string
		StrictHostKeyChecking bool
	}
	type args struct {
		args []string
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantErr      bool
		wantNotFound bool
	}{
		{
			name: "success",
			fields: fields{
				Config: &config.Project{EnvDir: func() string {
					temp, _ := os.MkdirTemp("", "test")
					return temp
				}()},
			},
			args: args{args: []string{"-M", "-t", "-S", "bastion.sock", "git@github.com", "-fN", "-o", "StrictHostKeyChecking=no"}},
		},
		{
			name: "not found",
			fields: fields{
				Config: &config.Project{EnvDir: func() string {
					temp, _ := os.MkdirTemp("", "test")
					return temp
				}()},
			},
			args:         args{args: []string{"-M", "-t", "-S", "bastion.sock", "git@github.com", "-fN", "-o", "StrictHostKeyChecking=no"}},
			wantErr:      false,
			wantNotFound: true,
		},
		{
			name: "incorrect host",
			fields: fields{
				Config: &config.Project{EnvDir: func() string {
					temp, _ := os.MkdirTemp("", "test")
					return temp
				}()},
			},
			args:         args{args: []string{"-M", "-t", "-S", "bastion.sock", "git@incorrecthost.com", "-fN", "-o", "StrictHostKeyChecking=no"}},
			wantErr:      true,
			wantNotFound: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &TunnelUpOptions{
				Config:                tt.fields.Config,
				PrivateKeyFile:        tt.fields.PrivateKeyFile,
				PublicKeyFile:         tt.fields.PublicKeyFile,
				BastionHostID:         tt.fields.BastionHostID,
				ForwardHost:           tt.fields.ForwardHost,
				StrictHostKeyChecking: tt.fields.StrictHostKeyChecking,
			}

			if err := o.runSSH(tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("runSSH() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantNotFound {
				tt.fields.Config.EnvDir = func() string {
					temp, _ := os.MkdirTemp("", "not_found")
					return temp
				}()
			}

			_, err := os.Stat(filepath.Join(tt.fields.Config.EnvDir, "bastion.sock"))
			if os.IsNotExist(err) != tt.wantNotFound {
				t.Error("bastion.sock not found")
				return
			}

			return
		})
	}
}

func Test_getPublicKey(t *testing.T) {
	tmp, _ := os.MkdirTemp("", "test")

	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "success", args: args{path: filepath.Join(tmp, "id_rsa.pub")}, want: func() string {
			pk, err := makeSSHKeyPair(filepath.Join(tmp, "id_rsa.pub"), filepath.Join(tmp, "id_rsa"))
			if err != nil {
				t.Fail()
			}
			return pk
		}(), wantErr: false},
		{name: "incorrect path", args: args{path: filepath.Join(tmp, "incorrect_path.pub")}, want: "", wantErr: true},
		{name: "invalid key", args: args{path: filepath.Join(tmp, "id_rsa.pub")}, want: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "invalid key" {
				err := ioutil.WriteFile(tt.args.path, []byte("invalid key"), os.ModeAppend)
				if err != nil {
					t.Fail()
				}
			}
			got, err := getPublicKey(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("getPublicKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getPublicKey() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func makeSSHKeyPair(pubKeyPath, privateKeyPath string) (string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return "", err
	}

	privateKeyFile, err := os.Create(privateKeyPath)
	defer func() {
		cerr := privateKeyFile.Close()
		if err == nil {
			err = cerr
		}
	}()
	if err != nil {
		return "", err
	}
	privateKeyPEM := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}
	if err := pem.Encode(privateKeyFile, privateKeyPEM); err != nil {
		return "", err
	}

	pub, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", err
	}

	var pubKeyBuf strings.Builder
	pubKeyBuf.Write(ssh.MarshalAuthorizedKey(pub))

	return pubKeyBuf.String(), ioutil.WriteFile(pubKeyPath, ssh.MarshalAuthorizedKey(pub), 0655)
}
