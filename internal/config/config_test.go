package config

import (
	"bytes"
	"fmt"
	"github.com/pterm/pterm"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/exec"
	"testing"
)

func TestSetTag(t *testing.T) {
	tests := []struct {
		name        string
		wantWarning bool
		tag         string
	}{
		{name: "success with git", tag: "git", wantWarning: false},
		{name: "success with ENV", tag: "env", wantWarning: true},
		{name: "with warning", wantWarning: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir, err := os.MkdirTemp("", "example")
			if err != nil {
				log.Fatal(err)
			}
			defer os.RemoveAll(dir) // clean up

			err = os.Chdir(dir)
			if err != nil {
				return
			}

			switch tt.tag {
			case "git":
				err = exec.Command("/bin/bash", "-c", "git init && touch file && git add file && git commit -m=\"test\"").Run()
				if err != nil {
					fmt.Println(err)
					return
				}
			case "env":
				viper.Set("ENV", "test")
			}

			buffer := &bytes.Buffer{}
			pterm.SetDefaultOutput(buffer)

			SetTag()

			fmt.Println(buffer.String())

			if buffer.String() != "" {
				if !tt.wantWarning {
					t.Fail()
				}
			}

			return
		})
	}
}
