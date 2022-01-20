package generate

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pterm/pterm"
)

const repoRegex = `((((git|hg)\+)?(git|ssh|file|https?):(//)?)|(\w+@[\w\.]+))`

func GenerateFiles(repoDir string, destionation string) (string, error) {
	return determineRepoDir(repoDir, destionation)
}

func determineRepoDir(template string, destionation string) (string, error) {
	if isRepoUrl(template) {
		return clone(template, destionation)
	} else {
		return "", fmt.Errorf("supported only repository url")
	}
}

func isRepoUrl(value string) bool {
	return regexp.MustCompile(repoRegex).Match([]byte(value))
}

func clone(url string, destination string) (string, error) {
	if destination == "" {
		destination = strings.Split(url, "/")[len(strings.Split(url, "/"))-1]
	}

	destination, err := filepath.Abs(destination)
	if err != nil {
		return "", fmt.Errorf("can't clone repository: %w", err)
	}

	cmd := exec.Command("git", "clone", url, destination)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("can't clone repository: %w: %s", err, stderr.String())
	}

	pterm.Info.Println(stderr.String())

	cmd = exec.Command("rm", "-rf", fmt.Sprintf("%s/.git", destination))

	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("can't clone repository: %w: %s", err, stderr.String())
	}

	return destination, nil
}
