package config

import (
	"bytes"
	"fmt"
	"github.com/pterm/pterm"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func Test_isStructured(t *testing.T) {
	tests := []struct {
		name        string
		dirName     string
		wantWarning bool
	}{
		{name: "success .infra", dirName: ".infra", wantWarning: false},
		{name: "success .ize", dirName: ".ize", wantWarning: false},
		{name: "with warning", dirName: ".test", wantWarning: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir, err := os.MkdirTemp("", "example")
			if err != nil {
				log.Fatal(err)
			}
			defer os.RemoveAll(dir) // clean up
			projectDir := filepath.Join(dir, tt.dirName)
			err = os.Mkdir(projectDir, defaultPerm)
			if err != nil {
				return
			}

			err = os.Chdir(dir)
			if err != nil {
				return
			}

			buffer := &bytes.Buffer{}

			pterm.SetDefaultOutput(buffer)
			isStructured()

			if buffer.String() == pterm.Warning.Sprint("is not an ize-structured directory. Please run ize init or cd into an ize-structured directory.\n") {
				if !tt.wantWarning {
					t.Fail()
				}
			} else {
				return
			}
		})
	}
}

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
