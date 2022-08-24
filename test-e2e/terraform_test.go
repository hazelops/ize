//go:build e2e && terraform
// +build e2e,terraform

package test

import (
	"fmt"
	"strings"
	"testing"
)

func TestIzeTerraformVersion(t *testing.T) {

	terraformVersionList := []string{
		"1.0.10",
		"1.1.3",
		"1.1.7",
		"1.2.6",
		"1.2.7",
	}

	if examplesRootDir == "" {
		t.Fatalf("Missing required environment variable IZE_EXAMPLES_PATH")
	}

	ize := NewBinary(t, izeBinary, examplesRootDir)

	for _, terraformVersion := range terraformVersionList {
		stdout, stderr, err := ize.RunRaw(fmt.Sprintf("--terraform-version=%s", terraformVersion), "terraform", "version")

		if err != nil {
			t.Errorf("error: %s", err)
		}

		if stderr != "" {
			t.Errorf("unexpected stderr output ize terraform version: %s", err)
		}

		if !strings.Contains(stdout, fmt.Sprintf("Terraform v%s", terraformVersion)) {
			t.Errorf("No success message detected after terraform version:\n%s", stdout)
		} else {
			t.Log(fmt.Sprintf("PASS: v%s: terraform version", terraformVersion))
		}
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
