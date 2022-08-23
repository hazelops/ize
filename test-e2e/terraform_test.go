//go:build e2e && terraform
// +build e2e,terraform

package test

import (
	"strings"
	"testing"
)

func TestIzeTerraformVersion(t *testing.T) {
	if examplesRootDir == "" {
		t.Fatalf("Missing required environment variable IZE_EXAMPLES_PATH")
	}

	ize := NewBinary(t, izeBinary, examplesRootDir)

	stdout, stderr, err := ize.RunRaw("terraform", "version")

	if err != nil {
		t.Errorf("error: %s", err)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output ize terraform version: %s", err)
	}

	if !strings.Contains(stdout, "Terraform v1.1.7") {
		t.Errorf("No success message detected after terraform version:\n%s", stdout)
	}
}

func TestIzeTerraformInit(t *testing.T) {
	if examplesRootDir == "" {
		t.Fatalf("Missing required environment variable IZE_EXAMPLES_PATH")
	}

	ize := NewBinary(t, izeBinary, examplesRootDir)

	stdout, stderr, err := ize.RunRaw("gen", "tfenv")

	if err != nil {
		t.Errorf("error: %s", err)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output ize gen tfenv: %s", err)
	}

	if !strings.Contains(stdout, "Generate terraform files completed") {
		t.Errorf("No success message detected after gen tfenv:\n%s", stdout)
	}

	stdout, stderr, err = ize.RunRaw("terraform", "init")

	if err != nil {
		t.Errorf("error: %s", err)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output ize terraform init: %s", err)
	}

	if !strings.Contains(stdout, "Terraform has been successfully initialized!") {
		t.Errorf("No success message detected after terraform init:\n%s", stdout)
	}
}
