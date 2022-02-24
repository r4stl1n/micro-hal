package drivers

import (
	"fmt"
	"time"

	i2c "github.com/r4stl1n/micro-hal/code/pkg/drivers/base"
)

const DefaultPCA9685Address = 0x40

// PCA9685 is a Driver for the PCA9685 16-channel 12-bit PWM/Servo controller
type PCA9685 struct {
	i2c     *i2c.I2C
	options *PCA9685Options
}

// PCA9685Options for controller
type PCA9685Options struct {
	Name       string
	StepCount  float32
	Frequency  float32
	ClockSpeed float32

	Mode1    byte
	PreScale byte
	Led0On   byte
}

// Init creates the new PCA9685 driver with specified i2c interface and options
func (pca9685 *PCA9685) Init(i2c *i2c.I2C, options *PCA9685Options) (*PCA9685, error) {

	adr := i2c.GetAddr()

	if i2c.GetAddr() == 0 {
		return nil, fmt.Errorf(`I2C device is not initiated`)
	}

	*pca9685 = PCA9685{
		i2c: i2c,
		options: &PCA9685Options{
			Name:       "PCA9685" + fmt.Sprintf("-0x%x", adr),
			StepCount:  4096.0,     // 12-bit
			Frequency:  50.0,       // 50Hz
			ClockSpeed: 25000000.0, // 25MHz

			Mode1:    0x00, // Default Mode1
			PreScale: 0xFE, // Default PreScale
			Led0On:   0x06, // default Led0On
		},
	}

	if options != nil {
		pca9685.options = options
	}

	if err := pca9685.i2c.WriteRegU8(pca9685.options.Mode1, 0x00|0xA1); err != nil { // Mode 1, autoincrement on)
		return nil, err
	}

	if err := pca9685.SetFreq(pca9685.options.Frequency); err != nil {
		return nil, err
	}

	return pca9685, nil
}

func (pca9685 *PCA9685) Name() string {
	return pca9685.options.Name
}

// SetFreq sets the PWM frequency in Hz for controller
func (pca9685 *PCA9685) SetFreq(freq float32) error {
	preScaleVal := pca9685.options.ClockSpeed/pca9685.options.StepCount/freq + 0.5

	if preScaleVal < 3.0 {
		return fmt.Errorf("PCA9685 cannot output at the given frequency")
	}

	oldMode, err := pca9685.i2c.ReadRegU8(pca9685.options.Mode1)

	if err != nil {
		return err
	}

	newMode := (oldMode & 0x7F) | 0x10 // Mode 1, sleep

	if err = pca9685.i2c.WriteRegU8(pca9685.options.Mode1, newMode); err != nil {
		return err
	}

	if err = pca9685.i2c.WriteRegU8(pca9685.options.PreScale, byte(preScaleVal)); err != nil {
		return err
	}

	if err = pca9685.i2c.WriteRegU8(pca9685.options.Mode1, oldMode); err != nil {
		return err
	}

	time.Sleep(5 * time.Millisecond)
	return nil
}

// GetFreq returns frequency value
func (pca9685 *PCA9685) GetFreq() float32 {
	return pca9685.options.Frequency
}

// Reset the chip
func (pca9685 *PCA9685) Reset() error {
	return pca9685.i2c.WriteRegU8(pca9685.options.Mode1, 0x00)
}

func (pca9685 *PCA9685) GetOptions() *PCA9685Options {
	return pca9685.options
}

// SetChannel sets a single PWM channel
func (pca9685 *PCA9685) SetChannel(chn, on, off int) error {

	if chn < 0 || chn > 15 {
		return fmt.Errorf("invalid [channel] value")
	}

	if on < 0 || on > int(pca9685.options.StepCount) {
		return fmt.Errorf("invalid [on] value")
	}

	if off < 0 || off > int(pca9685.options.StepCount) {
		return fmt.Errorf("invalid [off] value")
	}

	buf := []byte{pca9685.options.Led0On + byte(4*chn), byte(on) & 0xFF, byte(on >> 8), byte(off) & 0xFF, byte(off >> 8)}

	_, err := pca9685.i2c.WriteBytes(buf)

	return err
}
