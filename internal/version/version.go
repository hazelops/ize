package version

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hazelops/ize/internal/requirements"
	"log"
	"net/http"
	"runtime"

	"github.com/Masterminds/semver"
	"github.com/pterm/pterm"
)

var (
	GitCommit string
	Version   = "development"
)

func FullVersionNumber() string {
	var versionString bytes.Buffer

	if Version == "development" {
		return "development"
	}

	fmt.Fprintf(&versionString, "%s", Version)
	if GitCommit != "" {
		fmt.Fprintf(&versionString, " (%s)", GitCommit)
	}

	return versionString.String()
}

func CheckLatestRealese() {
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

	if Version != gr.Version {
		pterm.Warning.Printfln("Newest version is %s current version is %s. Consider upgrading.", gr.Version, Version)
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
