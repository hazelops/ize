package test

import (
	"io/fs"
	"path/filepath"
	"strings"
	"testing"
)

func TestIzeDeployAll(t *testing.T) {
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

	stdout, stderr, err := ize.RunRaw("deploy", "--auto-approve")

	if err != nil {
		t.Errorf("error: %s", err)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output ize deploy all: %s", err)
	}

	if !strings.Contains(stdout, "Deploy all completed!") {
		t.Errorf("No success message detected after all deploy:\n%s", stdout)
	}
}

func TestIzeDestroyAll(t *testing.T) {
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

	stdout, stderr, err := ize.RunRaw("destroy", "--auto-approve")

	if err != nil {
		t.Errorf("error: %s", err)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output ize deploy all: %s", err)
	}

	if !strings.Contains(stdout, "Destroy all completed!") {
		t.Errorf("No success message detected after all destroy:\n%s", stdout)
	}
}
