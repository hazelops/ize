//go:build e2e && ecs_cron
// +build e2e,ecs_cron

package test

import (
	"os"
	"strings"
	"testing"
)

// TestIzeCronUpInfra deploys the infrastructure (VPC, ECS cluster, cronjob task definition, EventBridge rule).
// This test reuses the ecs-apps-monorepo example which now includes a cronjob module.
func TestIzeCronUpInfra(t *testing.T) {
	if examplesRootDir == "" {
		t.Fatalf("Missing required environment variable IZE_EXAMPLES_PATH")
	}

	defer recovery(t)

	ize := NewBinary(t, izeBinary, examplesRootDir)

	stdout, stderr, err := ize.RunRaw("up", "infra")

	if err != nil {
		t.Errorf("error: %s", err)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output ize up infra: %s", stderr)
	}

	if !strings.Contains(stdout, "Deploy infra completed!") {
		t.Errorf("No success message detected after ize up infra:\n%s", stdout)
	}

	if os.Getenv("RUNNER_DEBUG") == "1" {
		t.Log(stdout)
	}
}

// TestIzeCronBuild builds the cronjob Docker image.
func TestIzeCronBuild(t *testing.T) {
	if examplesRootDir == "" {
		t.Fatalf("Missing required environment variable IZE_EXAMPLES_PATH")
	}

	defer recovery(t)

	ize := NewBinary(t, izeBinary, examplesRootDir)

	stdout, stderr, err := ize.RunRaw("build", "cronjob")

	if err != nil {
		t.Errorf("error: %s", err)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output ize build cronjob: %s", stderr)
	}

	if os.Getenv("RUNNER_DEBUG") == "1" {
		t.Log(stdout)
	}
}

// TestIzeCronPush pushes the cronjob image to ECR.
func TestIzeCronPush(t *testing.T) {
	if examplesRootDir == "" {
		t.Fatalf("Missing required environment variable IZE_EXAMPLES_PATH")
	}

	defer recovery(t)

	ize := NewBinary(t, izeBinary, examplesRootDir)

	stdout, stderr, err := ize.RunRaw("push", "cronjob")

	if err != nil {
		t.Errorf("error: %s", err)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output ize push cronjob: %s", stderr)
	}

	if os.Getenv("RUNNER_DEBUG") == "1" {
		t.Log(stdout)
	}
}

// TestIzeCronDeploy deploys the cron task — registers new task definition revision
// and updates the EventBridge target. Since no schedule is set in ize.toml,
// the existing Terraform-managed rule is kept as-is.
func TestIzeCronDeploy(t *testing.T) {
	if examplesRootDir == "" {
		t.Fatalf("Missing required environment variable IZE_EXAMPLES_PATH")
	}

	defer recovery(t)

	ize := NewBinary(t, izeBinary, examplesRootDir)

	stdout, stderr, err := ize.RunRaw("deploy", "cronjob")

	if err != nil {
		t.Errorf("error: %s", err)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output ize deploy cronjob: %s", stderr)
	}

	if !strings.Contains(stdout, "cron task deployed") {
		t.Errorf("No success message detected after ize deploy cronjob:\n%s", stdout)
	}

	if os.Getenv("RUNNER_DEBUG") == "1" {
		t.Log(stdout)
	}
}

// TestIzeCronRun manually triggers the cron task outside its schedule.
func TestIzeCronRun(t *testing.T) {
	if examplesRootDir == "" {
		t.Fatalf("Missing required environment variable IZE_EXAMPLES_PATH")
	}

	defer recovery(t)

	ize := NewBinary(t, izeBinary, examplesRootDir)

	stdout, stderr, err := ize.RunRaw("cron", "run", "cronjob")

	if err != nil {
		t.Errorf("error: %s", err)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output ize cron run cronjob: %s", stderr)
	}

	if !strings.Contains(stdout, "Task cronjob started") {
		t.Errorf("No success message detected after ize cron run cronjob:\n%s", stdout)
	}

	if os.Getenv("RUNNER_DEBUG") == "1" {
		t.Log(stdout)
	}
}

// TestIzeCronLogs fetches logs from the last cron task execution.
func TestIzeCronLogs(t *testing.T) {
	if examplesRootDir == "" {
		t.Fatalf("Missing required environment variable IZE_EXAMPLES_PATH")
	}

	defer recovery(t)

	ize := NewBinary(t, izeBinary, examplesRootDir)

	stdout, stderr, err := ize.RunRaw("cron", "logs", "cronjob")

	if err != nil {
		t.Errorf("error: %s", err)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output ize cron logs cronjob: %s", stderr)
	}

	// Logs may or may not be available depending on timing,
	// but the command should not error out
	if os.Getenv("RUNNER_DEBUG") == "1" {
		t.Log(stdout)
	}
}

// TestIzeCronDeployWithSchedule tests deploying with an explicit schedule override.
// This creates/updates the EventBridge rule with a new schedule.
func TestIzeCronDeployWithSchedule(t *testing.T) {
	if examplesRootDir == "" {
		t.Fatalf("Missing required environment variable IZE_EXAMPLES_PATH")
	}

	defer recovery(t)

	// TODO: When CLI flag --schedule is implemented, test schedule override here.
	// For now this test verifies the deploy path still works after the first deploy.

	ize := NewBinary(t, izeBinary, examplesRootDir)

	stdout, stderr, err := ize.RunRaw("deploy", "cronjob")

	if err != nil {
		t.Errorf("error: %s", err)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output ize deploy cronjob: %s", stderr)
	}

	if !strings.Contains(stdout, "cron task deployed") {
		t.Errorf("No success message detected after ize deploy cronjob:\n%s", stdout)
	}

	if os.Getenv("RUNNER_DEBUG") == "1" {
		t.Log(stdout)
	}
}

// TestIzeCronDown destroys the cron task (EventBridge rule + targets + task definitions).
func TestIzeCronDown(t *testing.T) {
	if examplesRootDir == "" {
		t.Fatalf("Missing required environment variable IZE_EXAMPLES_PATH")
	}

	defer recovery(t)

	ize := NewBinary(t, izeBinary, examplesRootDir)

	stdout, stderr, err := ize.RunRaw("down", "--auto-approve")

	if err != nil {
		t.Errorf("error: %s", err)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output ize down: %s", stderr)
	}

	if !strings.Contains(stdout, "Destroy all completed!") {
		t.Errorf("No success message detected after down:\n%s", stdout)
	}

	if os.Getenv("RUNNER_DEBUG") == "1" {
		t.Log(stdout)
	}
}
