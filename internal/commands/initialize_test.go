package commands

import (
	"testing"
)

func Test_writeConfig(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		exists  map[string]string
		wantErr bool
	}{
		{name: "success", path: "/tmp/ize.toml", exists: map[string]string{"namespace": "test"}, wantErr: false},
		{name: "invalid path", path: "/invalid/path/ize.toml", exists: map[string]string{}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := writeConfig(tt.path, tt.exists); (err != nil) != tt.wantErr {
				t.Errorf("writeConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
