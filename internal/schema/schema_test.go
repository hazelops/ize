package schema

import (
	"testing"
)

var valid = map[string]interface{}{
	"apps_path":   "/home/testnut/example/apps",
	"aws_profile": "testnut",
	"aws_region":  "us-east-1",
	"ecs": map[string]interface{}{
		"goblin":  map[string]interface{}{"cluster": "testnut-nutcorp", "skip_deploy": true, "timeout": 600},
		"squibby": map[string]interface{}{"timeout": 1200, "unsafe": true}},
	"env":            "testnut",
	"env_dir":        "/home/testnut/example/.ize/env/testnut",
	"home":           "/home/testnut",
	"ize_dir":        "/home/testnut/example/.ize",
	"namespace":      "testnut",
	"prefer_runtime": "native",
	"root_dir":       "/home/testnut/ize/example",
	"terraform": map[string]interface{}{
		"infra": map[string]interface{}{"aws_profile": "testnut", "root_domain_name": "examples.ize.sh", "version": "1.1.6"}},
	"terraform_version": "1.2.6",
	"tf_log":            "",
}

var invalidParameter = map[string]interface{}{
	"invalid":        "invalid",
	"aws_profile":    "testnut",
	"aws_region":     "us-east-1",
	"env":            "testnut",
	"env_dir":        "/home/testnut/example/.ize/env/testnut",
	"home":           "/home/testnut",
	"ize_dir":        "/home/testnut/example/.ize",
	"namespace":      "testnut",
	"prefer_runtime": "native",
	"root_dir":       "/home/testnut/ize/example",
}

var invalidType = map[string]interface{}{
	"aws_profile":    true,
	"aws_region":     "us-east-1",
	"env":            "testnut",
	"env_dir":        "/home/testnut/example/.ize/env/testnut",
	"home":           "/home/testnut",
	"ize_dir":        "/home/testnut/example/.ize",
	"namespace":      "testnut",
	"prefer_runtime": "native",
	"root_dir":       "/home/testnut/ize/example",
}

var validDeprecated = map[string]interface{}{
	"apps_path":   "/home/testnut/example/apps",
	"aws_profile": "testnut",
	"aws_region":  "us-east-1",
	"app": map[string]interface{}{
		"goblin":  map[string]interface{}{"cluster": "testnut-nutcorp", "skip_deploy": true, "timeout": 600},
		"squibby": map[string]interface{}{"timeout": 1200, "unsafe": true}},
	"env":            "testnut",
	"env_dir":        "/home/testnut/example/.ize/env/testnut",
	"home":           "/home/testnut",
	"ize_dir":        "/home/testnut/example/.ize",
	"namespace":      "testnut",
	"prefer_runtime": "native",
	"root_dir":       "/home/testnut/ize/example",
	"infra": map[string]interface{}{
		"terraform": map[string]interface{}{"aws_profile": "testnut", "root_domain_name": "examples.ize.sh", "version": "1.1.6"}},
	"terraform_version": "1.2.6",
	"tf_log":            "",
}

var empty = map[string]interface{}{}

func TestValidate(t *testing.T) {
	type args struct {
		config map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "valid", args: args{config: valid}, wantErr: false},
		{name: "valid deprecated", args: args{config: validDeprecated}, wantErr: false},
		{name: "invalid parameter", args: args{config: invalidParameter}, wantErr: true},
		{name: "invalid type", args: args{config: invalidType}, wantErr: true},
		{name: "empty", args: args{config: map[string]interface{}{}}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Validate(tt.args.config); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
