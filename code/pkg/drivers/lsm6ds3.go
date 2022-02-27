package drivers

import (
	"bytes"
	"encoding/binary"
	"fmt"

	i2c "github.com/r4stl1n/micro-hal/code/pkg/drivers/base"
)

const DefaultLSM6DS3Address = 0x6b

const (
	// lsm6ds3WhoAmI        = 0x0F
	// lsm6ds3Status        = 0x1E
	lsm6ds3Ctrl1XL = 0x10
	lsm6ds3Ctrl2G  = 0x11
	// lsm6ds3Ctrl3C         = 0x12
	lsm6ds3Ctrl4C = 0x13
	// lsm6ds3Ctrl5C            = 0x14
	// lsm6ds3Ctrl6C            = 0x15
	// lsm6ds3Ctrl7C            = 0x16
	// lsm6ds3Ctrl8Xl           = 0x17
	// lsm6ds3Ctrl9Xl           = 0x18
	// lsm6ds3Ctrl10C           = 0x19
	lsm6ds3OutXLG = 0x22
	// lsm6ds3OutXHG            = 0x23
	// lsm6ds3OutYLG            = 0x24
	// lsm6ds3OutYHG            = 0x25
	// lsm6ds3OutZLG            = 0x26
	// lsm6ds3OutZHG            = 0x27
	lsm6ds3OutLXL = 0x28
	// lsm6ds3OutXHXL           = 0x29
	// lsm6ds3OutYLXL           = 0x2A
	// lsm6ds3OutYHXL           = 0x2B
	// lsm6ds3OutZLXL           = 0x2C
	// lsm6ds3OutZHXL           = 0x2D
	lsm6ds3OutTempL = 0x20
	// lsm6ds3OutTempH          = 0x21
	// lsm6ds3BwScalOdrDisabled = 0x00
	lsm6ds3BwScalOdrEnabled = 0x80
	// lsm6ds3StepTimestampL    = 0x49
	// lsm6ds3StepTimestampH    = 0x4A
	// lsm6ds3StepCounterL      = 0x4B
	// lsm6ds3StepCounterH      = 0x4C
	// lsm6ds3StepCounterDelta  = 0x15
	// lsm6ds3TapCfg            = 0x58
	// lsm6ds3Int1Ctrl          = 0x0D

	lsm6ds3Accel2G uint8 = 0x00
	// lsm6ds3Accel4G  uint8 = 0x08
	// lsm6ds3Accel8G  uint8 = 0x0C
	// lsm6ds3Accel16G uint8 = 0x04

	// lsm6ds3AccelSrOFF   uint8 = 0x00
	// lsm6ds3AccelSr13    uint8 = 0x10
	// lsm6ds3AccelSr26    uint8 = 0x20
	// lsm6ds3AccelSr52    uint8 = 0x30
	lsm6ds3AccelSr104 uint8 = 0x40
	// lsm6ds3AccelSr208   uint8 = 0x50
	// lsm6ds3AccelSr416   uint8 = 0x60
	// lsm6ds3AccelSr833   uint8 = 0x70
	// lsm6ds3AccelSr1666  uint8 = 0x80
	// lsm6ds3AccelSr332   uint8 = 0x90
	// lsm6ds3AccelSr6664  uint8 = 0xA0
	// lsm6ds3AccelSr13330 uint8 = 0xB0

	// lsm6ds3AccelBw50  uint8 = 0x03
	lsm6ds3AccelBw100 uint8 = 0x02
	// lsm6ds3AccelBw200 uint8 = 0x01
	// lsm6ds3AccelBw400 uint8 = 0x00

	// lsm6ds3Gyro125Dps  uint8 = 0x01
	// lsm6ds3Gyro250Dps  uint8 = 0x00
	// lsm6ds3Gyro500Dps  uint8 = 0x04
	// lsm6ds3Gyro1000Dps uint8 = 0x08
	lsm6ds3Gyro2000Dps uint8 = 0x0C

	// lsm6ds3GyroSrOFF  uint8 = 0x00
	// lsm6ds3GyroSr13   uint8 = 0x10
	// lsm6ds3GyroSr26   uint8 = 0x20
	// lsm6ds3GyroSr52   uint8 = 0x30
	lsm6ds3GyroSr104 uint8 = 0x40
	// lsm6ds3GyroSr208  uint8 = 0x50
	// lsm6ds3GyroSr416  uint8 = 0x60
	// lsm6ds3GyroSr833  uint8 = 0x70
	// lsm6ds3GyroSr1666 uint8 = 0x80
)

type LSM6DS3Data struct {
	X int32
	Y int32
	Z int32
}

// LSM6DS3 is a Driver for the LSM6DS3 6-axis Accelerometer Gyroscope Sensor
type LSM6DS3 struct {
	i2c     *i2c.I2C
	options *LSM6DS3Options
}

// LSM6DS3Options for controller
type LSM6DS3Options struct {
	Name      string
	InCelsius bool
}

// Init creates the new LSM6DS3 driver with specified i2c interface and options
func (lsm6ds3 *LSM6DS3) Init(i2c *i2c.I2C, options *LSM6DS3Options) (*LSM6DS3, error) {

	adr := i2c.GetAddr()

	if i2c.GetAddr() == 0 {
		return nil, fmt.Errorf(`I2C device is not initiated`)
	}

	*lsm6ds3 = LSM6DS3{
		i2c: i2c,
		options: &LSM6DS3Options{
			Name:      "LSM6DS3" + fmt.Sprintf("-0x%x", adr),
			InCelsius: false,
		},
	}

	if options != nil {
		lsm6ds3.options = options
	}

	err := lsm6ds3.initProcess()

	if err != nil {
		return nil, err
	}

	return lsm6ds3, nil
}

func (lsm6ds3 *LSM6DS3) initProcess() error {

	// configure accelerometer mode
	_, err := lsm6ds3.i2c.WriteBytes([]byte{
		lsm6ds3Ctrl1XL,
		lsm6ds3Accel2G | lsm6ds3AccelSr104 | lsm6ds3AccelBw100,
	})

	if err != nil {
		return err
	}

	// Set ODR bit
	ctrl4cData, err := lsm6ds3.i2c.ReadRegU8(lsm6ds3Ctrl4C)

	if err != nil {
		return err
	}

	ctrl4cData = ctrl4cData &^ lsm6ds3BwScalOdrEnabled
	ctrl4cData |= lsm6ds3BwScalOdrEnabled

	_, err = lsm6ds3.i2c.WriteBytes([]byte{
		lsm6ds3Ctrl4C,
		ctrl4cData,
	})

	// Configure gyroscope
	_, err = lsm6ds3.i2c.WriteBytes([]byte{
		lsm6ds3Ctrl2G,
		lsm6ds3Gyro2000Dps | lsm6ds3GyroSr104,
	})

	if err != nil {
		return err
	}

	return err
}

func (lsm6ds3 *LSM6DS3) Name() string {
	return lsm6ds3.options.Name
}

func (lsm6ds3 *LSM6DS3) ReadData() (LSM6DS3Data, LSM6DS3Data, float32, error) {

	accelerationData, err := lsm6ds3.ReadAccelerationData()

	if err != nil {
		return LSM6DS3Data{}, LSM6DS3Data{}, 0, err
	}

	gyroscopeData, err := lsm6ds3.ReadGyroData()

	if err != nil {
		return LSM6DS3Data{}, LSM6DS3Data{}, 0, err
	}

	temperatureData, err := lsm6ds3.ReadTemperatureData()

	if err != nil {
		return LSM6DS3Data{}, LSM6DS3Data{}, 0, err
	}

	return accelerationData, gyroscopeData, temperatureData, nil
}

func (lsm6ds3 *LSM6DS3) ReadAccelerationData() (LSM6DS3Data, error) {

	data, _, err := lsm6ds3.i2c.ReadRegBytes(lsm6ds3OutLXL, 6)

	if err != nil {
		return LSM6DS3Data{}, err
	}

	x := int32(int16((uint16(data[1])<<8)|uint16(data[0]))) * 61
	y := int32(int16((uint16(data[3])<<8)|uint16(data[2]))) * 61
	z := int32(int16((uint16(data[5])<<8)|uint16(data[4]))) * 61

	// Convert from micro gravity to gravity values
	return LSM6DS3Data{X: x / 1000000, Y: y / 1000000, Z: z / 1000000}, nil
}

func (lsm6ds3 *LSM6DS3) ReadGyroData() (LSM6DS3Data, error) {

	data, _, err := lsm6ds3.i2c.ReadRegBytes(lsm6ds3OutXLG, 6)

	if err != nil {
		return LSM6DS3Data{}, err
	}

	x := int32(int16((uint16(data[1])<<8)|uint16(data[0]))) * 70000
	y := int32(int16((uint16(data[3])<<8)|uint16(data[2]))) * 70000
	z := int32(int16((uint16(data[5])<<8)|uint16(data[4]))) * 70000

	// Convert from micro degree to degrees
	return LSM6DS3Data{X: x / 1000000, Y: y / 1000000, Z: z / 1000000}, nil
}

func (lsm6ds3 *LSM6DS3) ReadTemperatureData() (float32, error) {

	data, _, err := lsm6ds3.i2c.ReadRegBytes(lsm6ds3OutTempL, 2)

	if err != nil {
		return 0.0, err
	}

	var temp int16

	buf := bytes.NewBuffer(data)

	err = binary.Read(buf, binary.BigEndian, &temp)

	// Convert the data to milli degree (C)
	t := float32(temp / 1000)

	if !lsm6ds3.options.InCelsius {
		// Convert to fahrenheit
		t = (t * 1.8) + 32
	}

	return t, nil

}
