//go:build e2e && terraform
// +build e2e,terraform

package test

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

// We're testing that we can download and run typical Terraform versions via ize
func TestIzeTerraformVersion_1_0_10(t *testing.T) {
	fmt.Println(os.Getenv("RUNNER_DEBUG"))

	terraformVersion := "1.0.10"
	if examplesRootDir == "" {
		t.Fatalf("Missing required environment variable IZE_EXAMPLES_PATH")
	}

	ize := NewBinary(t, izeBinary, examplesRootDir)

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
func TestIzeTerraformVersion_1_1_3(t *testing.T) {

	terraformVersion := "1.1.3"

	if examplesRootDir == "" {
		t.Fatalf("Missing required environment variable IZE_EXAMPLES_PATH")
	}

	ize := NewBinary(t, izeBinary, examplesRootDir)

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
func TestIzeTerraformVersion_1_1_7(t *testing.T) {

	terraformVersion := "1.1.7"

	if examplesRootDir == "" {
		t.Fatalf("Missing required environment variable IZE_EXAMPLES_PATH")
	}

	ize := NewBinary(t, izeBinary, examplesRootDir)

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
func TestIzeTerraformVersion_1_2_6(t *testing.T) {

	terraformVersion := "1.2.6"

	if examplesRootDir == "" {
		t.Fatalf("Missing required environment variable IZE_EXAMPLES_PATH")
	}

	ize := NewBinary(t, izeBinary, examplesRootDir)

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
func TestIzeTerraformVersion_1_2_7(t *testing.T) {

	terraformVersion := "1.2.7"

	if examplesRootDir == "" {
		t.Fatalf("Missing required environment variable IZE_EXAMPLES_PATH")
	}

	ize := NewBinary(t, izeBinary, examplesRootDir)

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

	if !strings.Contains(stdout, "Generate terraform file for \"infra\" completed") {
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
