//go:build e2e && ecs_apps
// +build e2e,ecs_apps

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

func TestIzeSecretsPushGoblin(t *testing.T) {
	if examplesRootDir == "" {
		t.Fatalf("Missing required environment variable IZE_EXAMPLES_PATH")
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

	if os.Getenv("RUNNER_DEBUG") == "1" {
		t.Log(stdout)
	}
}

func TestIzeSecretsPushSquibby(t *testing.T) {
	if examplesRootDir == "" {
		t.Fatalf("Missing required environment variable IZE_EXAMPLES_PATH")
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

	if os.Getenv("RUNNER_DEBUG") == "1" {
		t.Log(stdout)
	}
}

func TestIzeUpInfra(t *testing.T) {
	if examplesRootDir == "" {
		t.Fatalf("Missing required environment variable IZE_EXAMPLES_PATH")
	}

	defer recovery(t)

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

	stdout, stderr, err := ize.RunRaw("up", "infra")

	if err != nil {
		t.Errorf("error: %s", err)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output ize up all: %s", err)
	}

	if !strings.Contains(stdout, "Deploy infra completed!") {
		t.Errorf("No success message detected after ize up infra:\n%s", stdout)
	}

	if os.Getenv("RUNNER_DEBUG") == "1" {
		t.Log(stdout)
	}

	time.Sleep(time.Minute)
}

func TestIzeUpGoblin(t *testing.T) {
	if examplesRootDir == "" {
		t.Fatalf("Missing required environment variable IZE_EXAMPLES_PATH")
	}

	defer recovery(t)

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

	stdout, stderr, err := ize.RunRaw("up", "goblin")

	if err != nil {
		t.Errorf("error: %s", err)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output ize up all: %s", err)
	}

	if !strings.Contains(stdout, "Deploy app goblin completed") {
		t.Errorf("No success message detected after ize up squibby:\n%s", stdout)
	}

	if os.Getenv("RUNNER_DEBUG") == "1" {
		t.Log(stdout)
	}
}

func TestIzeUpSquibby(t *testing.T) {
	if examplesRootDir == "" {
		t.Fatalf("Missing required environment variable IZE_EXAMPLES_PATH")
	}

	defer recovery(t)

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

	stdout, stderr, err := ize.RunRaw("up", "squibby")

	if err != nil {
		t.Errorf("error: %s", err)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output ize up all: %s", err)
	}

	if !strings.Contains(stdout, "squibby") {
		t.Errorf("No success message detected after ize up squibby:\n%s", stdout)
	}

	if os.Getenv("RUNNER_DEBUG") == "1" {
		t.Log(stdout)
	}
}

func TestCheckSecretsSquibby(t *testing.T) {
	defer recovery(t)

	url := fmt.Sprintf("http://squibby.%s.examples.ize.sh/", os.Getenv("ENV"))

	for i := 0; i < 10; i++ {
		resp, err := http.Get(url)
		if err != nil {
			t.Error(err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
			t.Log(string(body))

		}

		if os.Getenv("RUNNER_DEBUG") == "1" {
			t.Logf("body: \n%s", string(body))
			t.Logf("status: \n%s", string(resp.Status))
		}

		if strings.Contains(string(body), exampleSquibbySecret) {
			return
		}

		time.Sleep(time.Second * 5)
	}

	t.Errorf("The expected string was not found in the response: %s", url)
}

func TestCheckSecretsGoblin(t *testing.T) {
	defer recovery(t)

	url := fmt.Sprintf("http://goblin.%s.examples.ize.sh/", os.Getenv("ENV"))

	for i := 0; i < 10; i++ {
		resp, err := http.Get(url)
		if err != nil {
			t.Error(err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
		}

		if os.Getenv("RUNNER_DEBUG") == "1" {
			t.Logf("body: \n%s", string(body))
			t.Logf("status: \n%s", string(resp.Status))
		}

		if strings.Contains(string(body), exampleGoblinSecret) {
			return
		}

		time.Sleep(time.Second * 5)
	}

	t.Errorf("The expected string was not found in the response: %s", url)
}

func TestIzeExecGoblin(t *testing.T) {
	if examplesRootDir == "" {
		t.Fatalf("Missing required environment variable IZE_EXAMPLES_PATH")
	}

	defer recovery(t)

	ize := NewBinary(t, izeBinary, examplesRootDir)

	stdout, stderr, err := ize.RunRaw("--plain-text-output", "exec", "goblin", "--", "sh -c \"echo $APP_NAME\"")

	if err != nil {
		t.Errorf("error: %s", err)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output ize exec goblin: %s", err)
	}

	if !strings.Contains(stdout, "goblin") || strings.Contains(stdout, "EOF") {
		t.Errorf("No success message detected after exec goblin:\n%s", stdout)
	}

	if os.Getenv("RUNNER_DEBUG") == "1" {
		t.Log(stdout)
	}
}

func TestIzeExecSquibby(t *testing.T) {
	if examplesRootDir == "" {
		t.Fatalf("Missing required environment variable IZE_EXAMPLES_PATH")
	}

	defer recovery(t)

	ize := NewBinary(t, izeBinary, examplesRootDir)

	stdout, stderr, err := ize.RunRaw("--plain-text-output", "exec", "squibby", "--", "sh -c \"echo $APP_NAME\"")

	if err != nil {
		t.Errorf("error: %s", err)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output ize exec squibby: %s", err)
	}

	if !strings.Contains(stdout, "squibby") || strings.Contains(stdout, "EOF") {
		t.Errorf("No success message detected after exec squibby:\n%s", stdout)
	}

	if os.Getenv("RUNNER_DEBUG") == "1" {
		t.Log(stdout)
	}
}

func TestIzeSecretsRmGoblin(t *testing.T) {
	if examplesRootDir == "" {
		t.Fatalf("Missing required environment variable IZE_EXAMPLES_PATH")
	}

	defer recovery(t)

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

	stdout, stderr, err := ize.RunRaw("secrets", "rm", "goblin")

	if err != nil {
		t.Errorf("error: %s", err)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output ize secret push: %s", err)
	}

	if !strings.Contains(stdout, "Removing secrets complete!") {
		t.Errorf("No success message detected after ize secret push:\n%s", stdout)
	}

	if os.Getenv("RUNNER_DEBUG") == "1" {
		t.Log(stdout)
	}
}

func TestIzeSecretsRmSquibby(t *testing.T) {
	if examplesRootDir == "" {
		t.Fatalf("Missing required environment variable IZE_EXAMPLES_PATH")
	}

	defer recovery(t)

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

	stdout, stderr, err := ize.RunRaw("secrets", "rm", "squibby")

	if err != nil {
		t.Errorf("error: %s", err)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr output ize secret push: %s", err)
	}

	if !strings.Contains(stdout, "Removing secrets complete!") {
		t.Errorf("No success message detected after ize secret push:\n%s", stdout)
	}

	if os.Getenv("RUNNER_DEBUG") == "1" {
		t.Log(stdout)
	}
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
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
		t.Errorf("unexpected stderr output ize dowb all: %s", err)
	}

	if !strings.Contains(stdout, "Destroy all completed!") {
		t.Errorf("No success message detected after down:\n%s", stdout)
	}

	if os.Getenv("RUNNER_DEBUG") == "1" {
		t.Log(stdout)
	}
}
