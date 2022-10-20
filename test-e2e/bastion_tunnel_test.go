//go:build e2e && bastion_tunnel
// +build e2e,bastion_tunnel

package test

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestIzeUpInfra(t *testing.T) {
	if examplesRootDir == "" {
		t.Fatalf("Missing required environment variable IZE_EXAMPLES_PATH")
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

	defer recovery(t)

	ize := NewBinary(t, izeBinary, examplesRootDir)

	stdout, stderr, err := ize.RunRaw("up", "infra")

	if err != nil {
		t.Errorf("error: %s", err)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output ize up all: %s", err)
	}

	if !strings.Contains(stdout, "Deploy infra completed!") {
		t.Errorf("No success message detected after all up:\n%s", stdout)
	}

	time.Sleep(time.Minute)

	t.Log(stdout)
}

func TestIzeTunnelUp(t *testing.T) {
	if examplesRootDir == "" {
		t.Fatalf("Missing required environment variable IZE_EXAMPLES_PATH")
	}

	defer recovery(t)

	ize := NewBinary(t, izeBinary, examplesRootDir)

	home, err := os.UserHomeDir()
	if err != nil {
		t.Errorf("error: %s", err)
	}

	time.Sleep(time.Minute)

	stdout, stderr, err := ize.RunRaw("tunnel", "up", "--ssh-public-key", filepath.Join(home, ".ssh", "id_rsa_tunnel_test.pub"), "--ssh-private-key", filepath.Join(home, ".ssh", "id_rsa_tunnel_test"))

	if err != nil {
		t.Errorf("error: %s", err)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output ize tunnel: %s", err)
	}

	t.Log(stdout)

	if !strings.Contains(stdout, "Tunnel is up! Forwarded ports:") {
		t.Errorf("No success message detected after tunnel:\n%s", stdout)
	}
}

func TestIzeTunnelStatus(t *testing.T) {
	if examplesRootDir == "" {
		t.Fatalf("Missing required environment variable IZE_EXAMPLES_PATH")
	}

	defer recovery(t)

	ize := NewBinary(t, izeBinary, examplesRootDir)

	stdout, stderr, err := ize.RunRaw("tunnel", "status")

	if err != nil {
		t.Errorf("error: %s", err)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output ize tunnel status: %s", err)
	}

	if !strings.Contains(stdout, "Tunnel is up. Forwarding config:") {
		t.Errorf("No success message detected after tunnel status:\n%s", stdout)
	}
}

func TestIzeTunnelDown(t *testing.T) {
	if examplesRootDir == "" {
		t.Fatalf("Missing required environment variable IZE_EXAMPLES_PATH")
	}

	defer recovery(t)

	ize := NewBinary(t, izeBinary, examplesRootDir)

	stdout, stderr, err := ize.RunRaw("tunnel", "down")

	if err != nil {
		t.Errorf("error: %s", err)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output ize tunnel down: %s", err)
	}

	if !strings.Contains(stdout, "Tunnel is down!") {
		t.Errorf("No success message detected after tunnel down:\n%s", stdout)
	}
}

func TestIzeDown(t *testing.T) {
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
		t.Errorf("unexpected stderr output ize down all: %s", err)
	}

	if !strings.Contains(stdout, "Destroy all completed!") {
		t.Errorf("No success message detected after all down:\n%s", stdout)
	}
}
