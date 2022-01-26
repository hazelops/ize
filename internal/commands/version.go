package commands

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Masterminds/semver"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var Version = "development"

func NewVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "show IZE version",
		Run: func(cmd *cobra.Command, args []string) {
			pterm.Info.Printfln("Version: %s", GetVersionNumber())
		},
	}

	return cmd
}

func GetVersionNumber() string {
	return Version
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
