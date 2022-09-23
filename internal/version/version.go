package version

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/hazelops/ize/internal/requirements"
	"github.com/pterm/pterm"
	"log"
	"net/http"
	"runtime"
	"time"
)

var (
	GitCommit string
	Version   = "development"
)

func FullVersionNumber() string {
	var versionString bytes.Buffer

	if Version == "development" {
		return fmt.Sprintf("development %s", time.Now().Format("2006-01-02T15:04:05"))
	}

	fmt.Fprintf(&versionString, "%s", Version)
	if GitCommit != "" {
		fmt.Fprintf(&versionString, " (%s)", GitCommit)
	}

	return versionString.String()
}

func CheckLatestRelease() {
	_, err := semver.NewVersion(Version)
	if err != nil {
		return
	}

	resp, err := http.Get("https://api.github.com/repos/hazelops/ize/releases/latest")
	if err != nil {
		log.Fatalln(err)
	}

	var gr gitResponse

	if err := json.NewDecoder(resp.Body).Decode(&gr); err != nil {
		log.Fatal(err)
	}

	var versionChangeAction = "upgrading"
	if Version > gr.Version {
		versionChangeAction = "downgrading"
	}
	if Version != gr.Version {
		pterm.Warning.Printfln("The newest stable version is %s, but your version is %s. Consider %s.", gr.Version, Version, versionChangeAction)
		ShowUpgradeCommand()
	}
}

type gitResponse struct {
	Version string `json:"tag_name"`
}

func ShowUpgradeCommand() error {
	switch goos := runtime.GOOS; goos {
	case "darwin":
		pterm.Warning.Println("Use the command to update\n:\tbrew upgrade ize")
	case "linux":
		distroName, err := requirements.ReadOSRelease("/etc/os-release")
		if err != nil {
			return err
		}
		switch distroName["ID"] {
		case "ubuntu":
			pterm.Warning.Println("Use the command to update:\n\tapt update && apt install ize")
		default:
			pterm.Warning.Println("See https://github.com/hazelops/ize/blob/main/DOCS.md#installation")
		}
	default:
		pterm.Warning.Println("See https://github.com/hazelops/ize/blob/main/DOCS.md#installation")
	}

	return nil
}
