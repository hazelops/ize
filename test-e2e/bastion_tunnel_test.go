//go:build e2e && bastion_tunnel
// +build e2e,bastion_tunnel

package test

import (
	"io/fs"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestIzeGenEnv_bastion_tunnel(t *testing.T) {
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
}

func TestIzeUpAll_bastion_tunnel(t *testing.T) {
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

	ize := NewBinary(t, izeBinary, examplesRootDir)

	stdout, stderr, err := ize.RunRaw("up", "--auto-approve")

	if err != nil {
		t.Errorf("error: %s", err)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output ize up all: %s", err)
	}

	if !strings.Contains(stdout, "Deploy all completed!") {
		t.Errorf("No success message detected after all up:\n%s", stdout)
	}

	time.Sleep(time.Minute)

	t.Log(stdout)
}

func TestIzeTunnelUp(t *testing.T) {
	if examplesRootDir == "" {
		t.Fatalf("Missing required environment variable IZE_EXAMPLES_PATH")
	}

	ize := NewBinary(t, izeBinary, examplesRootDir)

	stdout, stderr, err := ize.RunRaw("tunnel", "up")

	if err != nil {
		t.Errorf("error: %s", err)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output ize tunnel: %s", err)
	}

	if !strings.Contains(stdout, "Tunnel is up! Forwarded ports:") {
		t.Errorf("No success message detected after tunnel:\n%s", stdout)
	}
}

func TestIzeTunnelStatus(t *testing.T) {
	if examplesRootDir == "" {
		t.Fatalf("Missing required environment variable IZE_EXAMPLES_PATH")
	}

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

func TestIzeDownAll_bastion_tunnel(t *testing.T) {
	if examplesRootDir == "" {
		t.Fatalf("Missing required environment variable IZE_EXAMPLES_PATH")
	}

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
