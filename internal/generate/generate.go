package generate

import (
	"bytes"
	"embed"
	"fmt"
	pp "github.com/psihachina/path-parser"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/hazelops/ize/examples"
	"github.com/pterm/pterm"
)

const repoRegex = `((((git|hg)\+)?(git|ssh|file|https?):(//)?)|(\w+@[\w\.]+))`

func GenerateFiles(repoDir string, destionation string) (string, error) {
	return determineRepoDir(repoDir, destionation)
}

func GetDataFromFile(template string) ([]byte, error) {
	o := pp.ParsePath(template)

	switch o.Protocol {
	case "file":
		open, err := os.Open(o.Href)
		if err != nil {
			return nil, err
		}
		all, err := io.ReadAll(open)
		if err != nil {
			return nil, err
		}

		return all, nil
	//case "http", "https":
	//	parse, err := url.Parse(template)
	//	if err != nil {
	//		return nil, err
	//	}
	//	fmt.Println(parse.Path)
	//	re := regexp.MustCompile(`(\w+)[\/](.+)[:|\/](.+)`)
	//	all := re.FindStringSubmatch(parse.Path)
	//	fmt.Println(all[1:])
	//	return nil, nil
	//case "ssh":
	//	re := regexp.MustCompile(`(\w+)[\/](.+)[:|\/](.+)`)
	//	all := re.FindStringSubmatch(o.Pathname)
	//	fmt.Println(all[1:])
	//
	//	all[2] = strings.TrimSuffix(all[2], ".git")
	//
	//	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", all[1], all[2], all[3])
	//
	//	get, err := http.Get(url)
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	body, err := ioutil.ReadAll(get.Body)
	//
	//	var gitResp map[string]interface{}
	//	err = json.Unmarshal(body, &gitResp)
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	du := gitResp["download_url"].(string)
	//
	//	get, err = http.Get(du)
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	body, err = ioutil.ReadAll(get.Body)
	//
	//	return body, nil
	default:
		return nil, fmt.Errorf("can't get data from %s: type %s not supported", o.Href, o.Protocol)
	}
}

func determineRepoDir(template string, destination string) (string, error) {
	if isRepoUrl(template) {
		return clone(template, destination)
	} else if isInternalTemplate(template) {
		if destination == "" {
			destination = strings.Split(template, "/")[len(strings.Split(template, "/"))-1]
		}
		err := copyEmbedExamples(examples.Examples, template, destination)
		if err != nil {
			return "", err
		}
		return destination, nil
	} else {
		return "", fmt.Errorf("supported only repository url or internal examples")
	}
}

func isRepoUrl(value string) bool {
	return regexp.MustCompile(repoRegex).Match([]byte(value))
}

func isInternalTemplate(value string) bool {
	_, err := examples.Examples.ReadDir(value)
	return err == nil
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

func copyEmbedExamples(fsys embed.FS, sourceDir string, targetDir string) error {
	subdirs, err := fsys.ReadDir(sourceDir)
	if err != nil {
		return err
	}
	for _, d := range subdirs {
		sourcePath := path.Join(sourceDir, d.Name())
		if d.IsDir() {
			err = copyEmbedExamples(fsys, path.Join(sourceDir, d.Name()), path.Join(targetDir, d.Name()))
			if err != nil {
				return err
			}
		} else {
			localPath := filepath.Join(targetDir, d.Name())

			content, err := fsys.ReadFile(sourcePath)
			if err != nil {
				return err
			}
			err = os.MkdirAll(filepath.Dir(localPath), 0755)
			if err != nil {
				return err
			}
			err = os.WriteFile(localPath, content, 0755)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
