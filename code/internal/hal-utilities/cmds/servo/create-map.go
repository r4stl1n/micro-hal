package servos

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	components "github.com/r4stl1n/micro-hal/code/pkg/components"
	pca9685 "github.com/r4stl1n/micro-hal/code/pkg/drivers"
	i2c "github.com/r4stl1n/micro-hal/code/pkg/drivers/base"
	"github.com/r4stl1n/micro-hal/code/pkg/structs"
)

type CreateMap struct {
}

func (cmd *CreateMap) Init() *CreateMap {
	*cmd = CreateMap{}

	return cmd
}

func (cmd *CreateMap) Command() *cobra.Command {
	return &cobra.Command{
		Use:                   "createMap",
		Aliases:               []string{"cm"},
		Args:                  cobra.ExactArgs(6),
		ArgAliases:            []string{"i2c-address", "servoCount", "min_impusle", "max_impulse", "step", "outputFile"},
		DisableFlagsInUseLine: true,
		Short:                 "guided servo calibration map creation",
		Run:                   cmd.Run,
	}
}

func (cmd *CreateMap) getConveretedValues(args []string) (int, float32, float32, float32, error) {
	// Need to covert our current arguments into values
	servoCount, err := strconv.Atoi(args[1])

	if err != nil {
		return 0, 0.0, 0.0, 0, err
	}

	minImpulse, err := strconv.ParseFloat(args[2], 32)

	if err != nil {
		return 0, 0.0, 0.0, 0, err
	}

	maxImpulse, err := strconv.ParseFloat(args[3], 32)

	if err != nil {
		return 0, 0.0, 0.0, 0, err
	}

	step, err := strconv.ParseFloat(args[4], 32)

	if err != nil {
		return 0, 0.0, 0.0, 0, err
	}

	return servoCount, float32(minImpulse), float32(maxImpulse), float32(step), nil

}

func (cmd *CreateMap) minImpulseCalibrate(pca *pca9685.PCA9685, servoId int, actuationRange int, minImpulse float32, maxImpulse float32, step float32) float32 {

	logrus.Infof("Starting test for minImpulse level starting at: %f", minImpulse)

	// Create a new servo component
	servo := new(components.Servo).New(pca, servoId, &components.ServoOptions{
		ActuationRange: actuationRange,
		MinPulse:       minImpulse,
		MaxPulse:       maxImpulse,
	})

	logrus.Infof("Sending the servo move command for angle of 0 using impulse: %f", minImpulse)

	// Move the servo to a specific angle
	servo.Angle(0)

	scanner := bufio.NewScanner(os.Stdin)

	didAdjust := false

	for {

		fmt.Print("Did the servo move? (Y/N): ")
		scanner.Scan()

		text := scanner.Text()

		if len(text) != 0 {

			if strings.ToLower(text) == "n" {
				if didAdjust {
					minImpulse = minImpulse + step
				}
				break
			}

			didAdjust = true
			minImpulse = minImpulse - step

			servo = new(components.Servo).New(pca, servoId, &components.ServoOptions{
				ActuationRange: actuationRange,
				MinPulse:       minImpulse,
				MaxPulse:       maxImpulse,
			})

			logrus.Infof("Sending the servo move command for angle of 0 using impulse: %f", minImpulse)

			// Move the servo to a specific angle
			servo.Angle(0)

		} else {
			fmt.Println("Invalid selection please use Y or N")
		}
	}

	return minImpulse

}

func (cmd *CreateMap) maxImpulseCalibrate(pca *pca9685.PCA9685, servoId int, actuationRange int, minImpulse float32, maxImpulse float32, step float32) float32 {

	logrus.Infof("Starting test for maxImpulse level starting at: %f", maxImpulse)

	// Create a new servo component
	servo := new(components.Servo).New(pca, servoId, &components.ServoOptions{
		ActuationRange: actuationRange,
		MinPulse:       minImpulse,
		MaxPulse:       maxImpulse,
	})

	logrus.Infof("Sending the servo move command for angle of %d using impulse: %f", actuationRange, maxImpulse)

	// Move the servo to a specific angle
	servo.Angle(actuationRange)

	scanner := bufio.NewScanner(os.Stdin)

	didAdjust := false

	for {

		fmt.Print("Did the servo move? (Y/N): ")
		scanner.Scan()

		text := scanner.Text()

		if len(text) != 0 {

			if strings.ToLower(text) == "n" {
				if didAdjust {
					maxImpulse = maxImpulse - step
				}
				break
			}

			didAdjust = true
			maxImpulse = maxImpulse + step

			servo = new(components.Servo).New(pca, servoId, &components.ServoOptions{
				ActuationRange: actuationRange,
				MinPulse:       minImpulse,
				MaxPulse:       maxImpulse,
			})

			logrus.Infof("Sending the servo move command for angle of %d using impulse: %f", actuationRange, maxImpulse)

			// Move the servo to a specific angle
			servo.Angle(180)

		} else {
			fmt.Println("Invalid selection please use Y or N")
		}
	}

	return maxImpulse

}

func (cmd *CreateMap) moveToDefault(pca *pca9685.PCA9685, servoId int, actuationRange int, minImpulse float32, maxImpulse float32, angle int) {

	logrus.Infof("Moving servo %d to default position %d", servoId, angle)

	// Create a new servo component
	servo := new(components.Servo).New(pca, servoId, &components.ServoOptions{
		ActuationRange: actuationRange,
		MinPulse:       minImpulse,
		MaxPulse:       maxImpulse,
	})

	// Move the servo to a specific angle
	servo.Angle(angle)
}

func (cmd *CreateMap) Run(_ *cobra.Command, args []string) {

	servoCount, minImpulse, maxImpulse, step, err := cmd.getConveretedValues(args)

	if err != nil {
		logrus.Fatal(err)
	}

	// We create a connection to the i2c interface on the raspberry pi
	logrus.Infof("Attempting to connect to the i2c address: %s", args[0])
	i2c, err := i2c.New(pca9685.Address, args[0], i2c.DEFAULT_I2C_ADDRESS)
	if err != nil {
		logrus.Fatal(err)
	}

	// Next we create the needed driver to connect to the pca9685
	logrus.Info("Creating new connection to pca9685")
	pca, err := pca9685.New(i2c, nil)
	if err != nil {
		logrus.Fatal(err)
	}

	servoMap := structs.ServoCalibrationMap{}

	for i := 1; i != servoCount+1; i++ {

		servoId := i - 1

		logrus.Infof("Starting to tune servo in pin: %d", servoId)

		servoAlias := ""
		fmt.Print("Please enter an alias for the servo: ")
		fmt.Scanf("%s", &servoAlias)

		actuationRange := 180
		fmt.Print("Please enter the max actuation for the servo:")
		fmt.Scanf("%d", &actuationRange)

		defaultPosition := 90
		fmt.Print("Please enter the default position for the servo:")
		fmt.Scanf("%d", &defaultPosition)

		newMinImpulse := cmd.minImpulseCalibrate(pca, servoId, actuationRange, minImpulse, maxImpulse, step)

		newMaxImpulse := cmd.maxImpulseCalibrate(pca, servoId, actuationRange, newMinImpulse, maxImpulse, step)

		logrus.Infof("Min impulse is: %f, Max Impulse is: %f", newMinImpulse, newMaxImpulse)

		servoMap.Servos = append(servoMap.Servos, structs.ServoCalibrationItem{
			Alias:           servoAlias,
			PinId:           servoId,
			ActuationRange:  actuationRange,
			MinPulse:        newMinImpulse,
			MaxPulse:        newMaxImpulse,
			DefaultPosition: defaultPosition,
		})

		cmd.moveToDefault(pca, servoId, actuationRange, newMinImpulse, newMaxImpulse, defaultPosition)
	}

	marshaled, err := json.MarshalIndent(servoMap, "", " ")

	if err != nil {
		logrus.Fatal("Failed to create servo map from data: %s", err)
	}

	fmt.Println(string(marshaled))

	err = ioutil.WriteFile(args[5], marshaled, 0644)

	if err != nil {
		logrus.Fatal("Could not write file: %s", err.Error())
	}

	logrus.Infof("Servo map saved to: %s", args[5])

}
