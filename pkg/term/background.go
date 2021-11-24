package term

import (
	"log"
	"os/exec"
)

func (r Runner) BackgroundRun(name string, args []string) error {
	cmd := exec.Command(name, args...)
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
