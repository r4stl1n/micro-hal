package main

import (
	"fmt"
	"github.com/r4stl1n/micro-hal/code/internal/body-node/managers"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

// setupCloseHandler creates a 'listener' on a new goroutine which will notify the
// program if it receives an interrupt from the OS. We then handle this by calling
// our clean-up procedure and exiting the program.
func setupCloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		os.Exit(0)
	}()
}

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})

	logrus.SetLevel(logrus.InfoLevel)
}

func main() {
	setupCloseHandler()

	serviceManager := new(managers.NodeManager).Init()

	logrus.Info("controller-node started")

	serviceError := serviceManager.Process()

	if serviceError != nil {
		logrus.Fatal(serviceError)
	}
}
