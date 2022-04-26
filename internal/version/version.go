package version

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	}
}

type gitResponse struct {
	Version string `json:"tag_name"`
}
