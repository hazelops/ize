//go:build e2e
// +build e2e

package test

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

var (
	exampleGoblinSecret = ""
	exampleGoblinApiKey = ""

	exampleSquibbySecret = ""
	exampleSquibbyApiKey = ""
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

func TestIzeSecretsPushGoblin(t *testing.T) {
	if examplesRootDir == "" {
		t.Fatalf("Missing required environment variable IZE_PROJECT_TEMPLATE_PATH")
	}

	rand.Seed(time.Now().UTC().UnixNano())

	b := make([]byte, 12)

	for i := 0; i < 12; i++ {
		b[i] = byte(randInt(48, 127))
	}
	exampleGoblinSecret = string(b)

	for i := 0; i < 12; i++ {
		b[i] = byte(randInt(48, 127))
	}
	exampleGoblinApiKey = string(b)

	data := map[string]interface{}{
		"EXAMPLE_SECRET":  exampleGoblinSecret,
		"EXAMPLE_API_KEY": exampleGoblinApiKey,
	}

	jsonString, _ := json.Marshal(data)

	secretPath := filepath.Join(examplesRootDir, ".ize/env", os.Getenv("ENV"), "secrets/goblin.json")

	err := ioutil.WriteFile(secretPath, jsonString, os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}

	ize := NewBinary(t, izeBinary, examplesRootDir)

	stdout, stderr, err := ize.RunRaw("secrets", "push", "--force", "goblin")

	if err != nil {
		t.Errorf("error: %s", err)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output ize secret push: %s", err)
	}

	if !strings.Contains(stdout, "Pushing secrets complete!") {
		t.Errorf("No success message detected after ize secret push:\n%s", stdout)
	}
}

func TestIzeSecretsPushSquibby(t *testing.T) {
	if examplesRootDir == "" {
		t.Fatalf("Missing required environment variable IZE_PROJECT_TEMPLATE_PATH")
	}

	rand.Seed(time.Now().UTC().UnixNano())

	b := make([]byte, 12)

	for i := 0; i < 12; i++ {
		b[i] = byte(randInt(48, 127))
	}
	exampleSquibbySecret = string(b)

	for i := 0; i < 12; i++ {
		b[i] = byte(randInt(48, 127))
	}
	exampleSquibbyApiKey = string(b)

	data := map[string]interface{}{
		"EXAMPLE_SECRET":  exampleSquibbySecret,
		"EXAMPLE_API_KEY": exampleSquibbyApiKey,
	}

	secretPath := filepath.Join(examplesRootDir, ".ize/env", os.Getenv("ENV"), "secrets/squibby.json")

	jsonString, _ := json.Marshal(data)
	err := ioutil.WriteFile(secretPath, jsonString, os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}

	ize := NewBinary(t, izeBinary, examplesRootDir)

	stdout, stderr, err := ize.RunRaw("secrets", "push", "--force", "squibby")

	if err != nil {
		t.Errorf("error: %s", err)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output ize secret push: %s", err)
	}

	if !strings.Contains(stdout, "Pushing secrets complete!") {
		t.Errorf("No success message detected after ize secret push:\n%s", stdout)
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

func TestIzeExecGoblin(t *testing.T) {
	if examplesRootDir == "" {
		t.Fatalf("Missing required environment variable IZE_PROJECT_TEMPLATE_PATH")
	}

	ize := NewBinary(t, izeBinary, examplesRootDir)

	stdout, stderr, err := ize.RunRaw("exec", "goblin", "--", "sh -c \"echo $APP_NAME\"")

	if err != nil {
		t.Errorf("error: %s", err)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output ize deploy all: %s", err)
	}

	if !strings.Contains(stdout, "goblin") {
		t.Errorf("No success message detected after all deploy:\n%s", stdout)
	}
}

func TestCheckSecrets(t *testing.T) {
	resp, err := http.Get("http://squibby.testnut.examples.ize.sh/")
	if err != nil {
		t.Error(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}

	if !strings.Contains(string(body), exampleSquibbySecret) {
		t.Errorf("The installed env variable is not detected: %s", string(body))
	}

	resp, err = http.Get("http://goblin.testnut.examples.ize.sh/")
	if err != nil {
		t.Error(err)
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}

	if !strings.Contains(string(body), exampleGoblinSecret) {
		t.Errorf("The installed env variable is not detected: %s", string(body))
	}
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
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
