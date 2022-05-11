package test

import (
	"io/fs"
	"path/filepath"
	"strings"
	"testing"
)

func TestIzeGenEnv(t *testing.T) {
	if examplesRootDir == "" {
		t.Fatalf("Missing required environment variable IZE_PROJECT_TEMPLATE_PATH")
	}

	ize := NewBinary(t, izeBinary, examplesRootDir)

	stdout, stderr, err := ize.RunRaw("gen", "env")

	if err != nil {
		t.Errorf("error: %s", err)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output ize deploy all: %s", err)
	}

	if !strings.Contains(stdout, "Generate terraform files completed") {
		t.Errorf("No success message detected after all deploy:\n%s", stdout)
	}
}

func TestIzeDeployAll(t *testing.T) {
	if examplesRootDir == "" {
		t.Fatalf("Missing required environment variable IZE_PROJECT_TEMPLATE_PATH")
	}

	foundIZEConfig := false
	err := filepath.Walk(examplesRootDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Name() == "ize.toml" {
			foundIZEConfig = true
		}

		return nil
	})
	if err != nil {
		t.Fatalf("Failed listing files in project template path %s: %s", examplesRootDir, err)
	}

	if !foundIZEConfig {
		t.Fatalf("No ize.toml file in project template path %s", examplesRootDir)
	}

	ize := NewBinary(t, izeBinary, examplesRootDir)

	stdout, stderr, err := ize.RunRaw("deploy", "--auto-approve", "--prefer-runtime=docker")

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
	if examplesRootDir == "" {
		t.Fatalf("Missing required environment variable IZE_PROJECT_TEMPLATE_PATH")
	}

	ize := NewBinary(t, izeBinary, examplesRootDir)

	stdout, stderr, err := ize.RunRaw("destroy", "--auto-approve", "--prefer-runtime=docker")

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
