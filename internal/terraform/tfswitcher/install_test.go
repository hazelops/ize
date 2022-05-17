//go:build !e2e
// +build !e2e

package tfswitcher

import (
	"os/user"
	"runtime"
	"testing"
)

func getInstallFile(installFile string) string {
	if runtime.GOOS == "windows" {
		return installFile + ".exe"
	}

	return installFile
}

func TestInstall(t *testing.T) {

	t.Run("User should exist",
		func(t *testing.T) {
			usr, errCurr := user.Current()
			if errCurr != nil {
				t.Errorf("Unable to get user %v [unexpected]", errCurr)
			}

			if usr != nil {
				t.Logf("Current user exist: %v  [expected]\n", usr.HomeDir)
			} else {
				t.Error("Unable to get user [unexpected]")
			}
		},
	)
}
