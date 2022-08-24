package version

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hazelops/ize/internal/config"
	"log"
	"net/http"

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
		config.ShowUpgradeCommand()
	}
}

type gitResponse struct {
	Version string `json:"tag_name"`
}
