package utils

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	drivers "github.com/r4stl1n/micro-hal/code/pkg/drivers"
	base "github.com/r4stl1n/micro-hal/code/pkg/drivers/base"
)

type LSM6DS3Test struct {
}

func (cmd *LSM6DS3Test) Init() *LSM6DS3Test {
	*cmd = LSM6DS3Test{}

	return cmd
}

func (cmd *LSM6DS3Test) Command() *cobra.Command {
	return &cobra.Command{
		Use:                   "lsm6ds3",
		Aliases:               []string{"lsm"},
		Args:                  cobra.ExactArgs(1),
		ArgAliases:            []string{"i2c-address"},
		DisableFlagsInUseLine: true,
		Short:                 "test lsm",
		Run:                   cmd.Run,
	}
}

func (cmd *LSM6DS3Test) Run(_ *cobra.Command, args []string) {

	// We create a connection to the i2c interface on the raspberry pi
	logrus.Infof("Attempting to connect to the i2c address: %s", args[0])
	i2c, err := new(base.I2C).Init(drivers.DefaultLSM6DS3Address, args[0], base.DEFAULT_I2C_ADDRESS)

	if err != nil {
		logrus.Fatal(err)
	}

	// Next we create the needed driver to connect to the pca9685
	logrus.Info("Creating new connection to LSM6DS3Test")

	lsm, err := new(drivers.LSM6DS3).Init(i2c, nil)

	if err != nil {
		logrus.Fatal(err)
	}

	for {
		accelerometer, gyroscope, temperature, err := lsm.ReadData()

		if err != nil {
			logrus.Fatal(err)
		}

		logrus.Infof("Accelerometer: %+v", accelerometer)
		logrus.Infof("Gyroscope: %+v", gyroscope)
		logrus.Infof("Temprature: %+v", temperature)

		time.Sleep(time.Duration(1) * time.Second)
	}
}
