package drivers

import (
	"fmt"

	i2c "github.com/r4stl1n/micro-hal/code/pkg/drivers/base"
)

const DefaultLSM6DS3Address = 0x6b

const (
	WHO_AM_I             = 0x0F
	STATUS               = 0x1E
	CTRL1_XL             = 0x10
	CTRL2_G              = 0x11
	CTRL3_C              = 0x12
	CTRL4_C              = 0x13
	CTRL5_C              = 0x14
	CTRL6_C              = 0x15
	CTRL7_G              = 0x16
	CTRL8_XL             = 0x17
	CTRL9_XL             = 0x18
	CTRL10_C             = 0x19
	OUTX_L_G             = 0x22
	OUTX_H_G             = 0x23
	OUTY_L_G             = 0x24
	OUTY_H_G             = 0x25
	OUTZ_L_G             = 0x26
	OUTZ_H_G             = 0x27
	OUTX_L_XL            = 0x28
	OUTX_H_XL            = 0x29
	OUTY_L_XL            = 0x2A
	OUTY_H_XL            = 0x2B
	OUTZ_L_XL            = 0x2C
	OUTZ_H_XL            = 0x2D
	OUT_TEMP_L           = 0x20
	OUT_TEMP_H           = 0x21
	BW_SCAL_ODR_DISABLED = 0x00
	BW_SCAL_ODR_ENABLED  = 0x80
	STEP_TIMESTAMP_L     = 0x49
	STEP_TIMESTAMP_H     = 0x4A
	STEP_COUNTER_L       = 0x4B
	STEP_COUNTER_H       = 0x4C
	STEP_COUNT_DELTA     = 0x15
	TAP_CFG              = 0x58
	INT1_CTRL            = 0x0D

	ACCEL_2G  uint8 = 0x00
	ACCEL_4G  uint8 = 0x08
	ACCEL_8G  uint8 = 0x0C
	ACCEL_16G uint8 = 0x04

	ACCEL_SR_OFF   uint8 = 0x00
	ACCEL_SR_13    uint8 = 0x10
	ACCEL_SR_26    uint8 = 0x20
	ACCEL_SR_52    uint8 = 0x30
	ACCEL_SR_104   uint8 = 0x40
	ACCEL_SR_208   uint8 = 0x50
	ACCEL_SR_416   uint8 = 0x60
	ACCEL_SR_833   uint8 = 0x70
	ACCEL_SR_1666  uint8 = 0x80
	ACCEL_SR_3332  uint8 = 0x90
	ACCEL_SR_6664  uint8 = 0xA0
	ACCEL_SR_13330 uint8 = 0xB0

	ACCEL_BW_50  uint8 = 0x03
	ACCEL_BW_100 uint8 = 0x02
	ACCEL_BW_200 uint8 = 0x01
	ACCEL_BW_400 uint8 = 0x00

	//GYRO_125DPS  uint8 = 0x01
	GYRO_250DPS  uint8 = 0x00
	GYRO_500DPS  uint8 = 0x04
	GYRO_1000DPS uint8 = 0x08
	GYRO_2000DPS uint8 = 0x0C

	GYRO_SR_OFF  uint8 = 0x00
	GYRO_SR_13   uint8 = 0x10
	GYRO_SR_26   uint8 = 0x20
	GYRO_SR_52   uint8 = 0x30
	GYRO_SR_104  uint8 = 0x40
	GYRO_SR_208  uint8 = 0x50
	GYRO_SR_416  uint8 = 0x60
	GYRO_SR_833  uint8 = 0x70
	GYRO_SR_1666 uint8 = 0x80
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

	//err := lsm6ds3.initProcess()

	//if err != nil {
	//		return nil, err
	//	}

	return lsm6ds3, nil
}

func (lsm6ds3 *LSM6DS3) initProcess() error {

	// configure accelerometer mode
	_, err := lsm6ds3.i2c.WriteBytes([]byte{
		CTRL1_XL,
		ACCEL_2G | ACCEL_SR_104 | ACCEL_BW_100,
	})

	if err != nil {
		return err
	}

	// Set ODR bit
	ctrl4cData, err := lsm6ds3.i2c.ReadRegU8(CTRL4_C)

	if err != nil {
		return err
	}

	ctrl4cData = ctrl4cData &^ BW_SCAL_ODR_ENABLED
	ctrl4cData |= BW_SCAL_ODR_ENABLED

	_, err = lsm6ds3.i2c.WriteBytes([]byte{
		CTRL4_C,
		ctrl4cData,
	})

	// Configure gyroscope
	_, err = lsm6ds3.i2c.WriteBytes([]byte{
		CTRL2_G,
		GYRO_2000DPS | GYRO_SR_104,
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

	data, _, err := lsm6ds3.i2c.ReadRegBytes(OUTX_L_XL, 6)

	if err != nil {
		return LSM6DS3Data{}, err
	}

	x := int32(int16((uint16(data[1])<<8)|uint16(data[0]))) * 61
	y := int32(int16((uint16(data[3])<<8)|uint16(data[2]))) * 61
	z := int32(int16((uint16(data[5])<<8)|uint16(data[4]))) * 61

	return LSM6DS3Data{X: x, Y: y, Z: z}, nil
}

func (lsm6ds3 *LSM6DS3) ReadGyroData() (LSM6DS3Data, error) {

	data, _, err := lsm6ds3.i2c.ReadRegBytes(OUTX_L_G, 6)

	if err != nil {
		return LSM6DS3Data{}, err
	}

	x := int32(int16((uint16(data[1])<<8)|uint16(data[0]))) * 70000
	y := int32(int16((uint16(data[3])<<8)|uint16(data[2]))) * 70000
	z := int32(int16((uint16(data[5])<<8)|uint16(data[4]))) * 70000

	return LSM6DS3Data{X: x, Y: y, Z: z}, nil
}

func (lsm6ds3 *LSM6DS3) ReadTemperatureData() (float32, error) {

	data, _, err := lsm6ds3.i2c.ReadRegBytes(OUT_TEMP_L, 62)

	if err != nil {
		return 0.0, err
	}

	// Convert the data to milli degree (C)
	temperature := float32(25000 + (int32((int16(data[1])<<8)|int16(data[0]))*125)/2)

	// Convert the milli degree to degree (C)
	temperature = temperature / 1000

	if !lsm6ds3.options.InCelsius {
		// Convert to fahrenheit
		temperature = (temperature * 1.8) + 32

	}

	return temperature, nil

}
