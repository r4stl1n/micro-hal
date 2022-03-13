package servos

import (
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	components "github.com/r4stl1n/micro-hal/code/pkg/components"
	drivers "github.com/r4stl1n/micro-hal/code/pkg/drivers"
	base "github.com/r4stl1n/micro-hal/code/pkg/drivers/base"
)

type Move struct {
}

func (cmd *Move) Init() *Move {
	*cmd = Move{}

	return cmd
}

func (cmd *Move) Command() *cobra.Command {
	return &cobra.Command{
		Use:                   "move",
		Aliases:               []string{"m"},
		Args:                  cobra.ExactArgs(6),
		ArgAliases:            []string{"i2c-address", "servo_id", "actuationRange", "min_impusle", "max_impulse", "angle"},
		DisableFlagsInUseLine: true,
		Short:                 "move servo",
		Run:                   cmd.Run,
	}
}

func (cmd *Move) getConveretedValues(args []string) (int, int, float32, float32, int, error) {
	// Need to covert our current arguments into values
	servoId, err := strconv.Atoi(args[1])

	if err != nil {
		return 0, 0, 0, 0, 0.0, err
	}

	actuationRange, err := strconv.Atoi(args[2])

	if err != nil {
		return 0, 0, 0, 0, 0.0, err
	}

	minImpulse, err := strconv.ParseFloat(args[3], 32)

	if err != nil {
		return 0, 0, 0, 0, 0.0, err
	}

	maxImpulse, err := strconv.ParseFloat(args[4], 32)

	if err != nil {
		return 0, 0, 0, 0, 0.0, err
	}

	angle, err := strconv.Atoi(args[5])

	if err != nil {
		return 0, 0, 0, 0, 0.0, err
	}

	return servoId, actuationRange, float32(minImpulse), float32(maxImpulse), angle, nil

}

func (cmd *Move) Run(_ *cobra.Command, args []string) {

	servoId, actuationRange, minImpulse, maxImpulse, angle, err := cmd.getConveretedValues(args)

	if err != nil {
		logrus.Fatal(err)
	}

	// We create a connection to the i2c interface on the raspberry pi
	logrus.Infof("Attempting to connect to the i2c address: %s", args[0])
	i2c, err := new(base.I2C).Init(drivers.DefaultPCA9685Address, args[0], base.DEFAULT_I2C_ADDRESS)
	if err != nil {
		logrus.Fatal(err)
	}

	// Next we create the needed driver to connect to the pca9685
	logrus.Info("Creating new connection to pca9685")
	pca, err := new(drivers.PCA9685).Init(i2c, nil)
	if err != nil {
		logrus.Fatal(err)
	}

	// Create a new servo component
	servo := new(components.Servo).Init(pca, servoId, &components.ServoOptions{
		ActuationRange: actuationRange,
		MinPulse:       minImpulse,
		MaxPulse:       maxImpulse,
	})

	logrus.Infof("Sending the servo move command for angle: %d", angle)

	// Move the servo to a specific angle
	servo.Angle(angle)

}
