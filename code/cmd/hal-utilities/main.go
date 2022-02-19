package main

import (
	"fmt"

	cli "github.com/r4stl1n/micro-hal/code/internal/hal-utilities"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})

	logrus.SetLevel(logrus.InfoLevel)
}

func main() {
	runError := new(cli.CLI).Init().Run().Error
	if runError != nil {
		fmt.Println(runError.Error())
	}
}
