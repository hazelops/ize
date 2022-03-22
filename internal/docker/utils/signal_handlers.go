package utils

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
)

func SetupSignalHandlers(cli *client.Client, containerID string) {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)

	go func() {
		for {
			select {
			case s := <-signalChannel:
				logrus.Debug("Received signal:", s)

				cli.ContainerKill(context.Background(), containerID, strconv.Itoa(int(s.(syscall.Signal))))
			}
		}
	}()
}
