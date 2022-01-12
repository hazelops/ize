package tunnel

import (
	"context"
	"fmt"
	"os"
	"syscall"

	"github.com/pterm/pterm"
	"github.com/sevlyar/go-daemon"
)

func daemonContext(c context.Context) *daemon.Context {
	return &daemon.Context{
		PidFileName: "tunnel.pid",
		PidFilePerm: 0644,
		LogFileName: "tunnel.log",
		LogFilePerm: 0640,
	}
}

func killDaemon(dCtx *daemon.Context) error {
	p, err := dCtx.Search()
	if err != nil {
		return fmt.Errorf("search for daemon process: %w", err)
	}
	pterm.Info.Printf("killing daemon process(pid: %d)\n", p.Pid)

	if err := p.Signal(syscall.SIGTERM); err != nil {
		return fmt.Errorf("kill daemon process(pid: %d): %w", p.Pid, err)
	}
	return os.Remove(dCtx.PidFileName)
}

func daemonRunning(dCtx *daemon.Context) (process *os.Process, running bool, err error) {
	p, err := dCtx.Search()
	if err != nil {
		if os.IsNotExist(err) {
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("search daemon process: %w", err)
	}
	err = p.Signal(syscall.Signal(0))
	if err != nil {
		return p, false, nil
	}
	return p, true, nil
}
