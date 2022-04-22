package test

import (
	"io/fs"
	"path/filepath"
	"strings"
	"testing"
)

func TestECS(t *testing.T) {
	izeBinary := GetFromEnv("IZE_BINARY", "ize")
	projectTemplatePath := GetFromEnv("IZE_PROJECT_TEMPLATE_PATH", "")

	if projectTemplatePath == "" {
		t.Fatalf("Missing required environment variable IZE_PROJECT_TEMPLATE_PATH")
	}

	foundIZEConfig := false
	err := filepath.Walk(projectTemplatePath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Name() == "ize.toml" {
			foundIZEConfig = true
		}

		return nil
	})
	if err != nil {
		t.Fatalf("Failed listing files in project template path %s: %s", projectTemplatePath, err)
	}

	if !foundIZEConfig {
		t.Fatalf("No ize.toml file in project template path %s", projectTemplatePath)
	}

	ize := NewBinary(t, izeBinary, projectTemplatePath)

	stdout, stderr, err := ize.RunRaw("deploy", "sqiubby")

	if err != nil {
		t.Errorf("error: %s", err)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output ize deploy squibby app: %s", err)
	}

	if !strings.Contains(stdout, "deploy app sqiubby completed") {
		t.Errorf("No success message detected after app deploy:\n%s", stdout)
	}
}

func TestServerless(t *testing.T) {
	izeBinary := GetFromEnv("IZE_BINARY", "ize")
	projectTemplatePath := GetFromEnv("IZE_PROJECT_TEMPLATE_PATH", "")

	if projectTemplatePath == "" {
		t.Fatalf("Missing required environment variable IZE_PROJECT_TEMPLATE_PATH")
	}

	foundIZEConfig := false
	err := filepath.Walk(projectTemplatePath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Name() == "ize.toml" {
			foundIZEConfig = true
		}

		return nil
	})
	if err != nil {
		t.Fatalf("Failed listing files in project template path %s: %s", projectTemplatePath, err)
	}

	if !foundIZEConfig {
		t.Fatalf("No ize.toml file in project template path %s", projectTemplatePath)
	}

	ize := NewBinary(t, izeBinary, projectTemplatePath)

	stdout, stderr, err := ize.RunRaw("deploy", "pecan")

	if err != nil {
		t.Errorf("error: %s", err)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output ize deploy squibby app: %s", err)
	}

	if !strings.Contains(stdout, "deploy app pecan completed") {
		t.Errorf("No success message detected after app deploy:\n%s", stdout)
	}
}
