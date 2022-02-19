package pca9685

import (
	"fmt"
	"time"

	i2c "github.com/r4stl1n/micro-hal/code/pkg/drivers/base"
)

const (
	// Address default for controller
	Address byte = 0x40

	// Registers
	Mode1    byte = 0x00
	Prescale byte = 0xFE
	Led0On   byte = 0x06

	// The internal reference clock is 25mhz but may vary slightly with
	// environmental conditions and manufacturing variances. Providing a more precise
	// "ReferenceClockSpeed" can improve the accuracy of the frequency and duty_cycle computations.
	ReferenceClockSpeed float32 = 25000000.0 // 25MHz
	StepCount           float32 = 4096.0     // 12-bit
	DefaultPWMFrequency float32 = 50.0       // 50Hz
)

// PCA9685 is a Driver for the PCA9685 16-channel 12-bit PWM/Servo controller
type PCA9685 struct {
	i2c  *i2c.Options
	optn *Options
}

// Options for controller
type Options struct {
	Name       string
	Frequency  float32
	ClockSpeed float32
}

// New creates the new PCA9685 driver with specified i2c interface and options
func New(i2c *i2c.Options, optn *Options) (*PCA9685, error) {
	adr := i2c.GetAddr()
	if adr == 0 {
		return nil, fmt.Errorf(`I2C device is not initiated`)
	}

	pca := &PCA9685{
		i2c: i2c,
		optn: &Options{
			Name:       "Controller" + fmt.Sprintf("-0x%x", adr),
			Frequency:  DefaultPWMFrequency,
			ClockSpeed: ReferenceClockSpeed,
		},
	}
	if optn != nil {
		pca.optn = optn
	}

	if err := pca.i2c.WriteRegU8(Mode1, 0x00|0xA1); err != nil { // Mode 1, autoincrement on)
		return nil, err
	}
	if err := pca.SetFreq(DefaultPWMFrequency); err != nil {
		return nil, err
	}
	return pca, nil
}

// SetFreq sets the PWM frequency in Hz for controller
func (pca *PCA9685) SetFreq(freq float32) (err error) {
	prescaleVal := ReferenceClockSpeed/StepCount/freq + 0.5
	if prescaleVal < 3.0 {
		return fmt.Errorf("PCA9685 cannot output at the given frequency")
	}
	oldMode, err := pca.i2c.ReadRegU8(Mode1)
	if err != nil {
		return err
	}
	newMode := (oldMode & 0x7F) | 0x10 // Mode 1, sleep
	if err := pca.i2c.WriteRegU8(Mode1, newMode); err != nil {
		return err
	}
	if err := pca.i2c.WriteRegU8(Prescale, byte(prescaleVal)); err != nil {
		return err
	}
	if err := pca.i2c.WriteRegU8(Mode1, oldMode); err != nil {
		return err
	}
	time.Sleep(5 * time.Millisecond)
	return nil
}

// GetFreq returns frequency value
func (pca *PCA9685) GetFreq() float32 {
	return pca.optn.Frequency
}

// Reset the chip
func (pca *PCA9685) Reset() (err error) {
	return pca.i2c.WriteRegU8(Mode1, 0x00)
}

// SetChannel sets a single PWM channel
func (pca *PCA9685) SetChannel(chn, on, off int) (err error) {
	if chn < 0 || chn > 15 {
		return fmt.Errorf("invalid [channel] value")
	}
	if on < 0 || on > int(StepCount) {
		return fmt.Errorf("invalid [on] value")
	}
	if off < 0 || off > int(StepCount) {
		return fmt.Errorf("invalid [off] value")
	}

	buf := []byte{Led0On + byte(4*chn), byte(on) & 0xFF, byte(on >> 8), byte(off) & 0xFF, byte(off >> 8)}
	_, err = pca.i2c.WriteBytes(buf)
	return err
}

const (
	// The specified pulse width range of a servo has historically been 1000-2000us,
	// for a 90 degree range of motion. But nearly all modern servos have a 170-180
	// degree range, and the pulse widths can go well out of the range to achieve this
	// extended motion. The default values here of `750` and `2250` typically give
	// 135 degrees of motion. You can set `Range` to correspond to the
	// actual range of motion you observe with your given `MinPulse` and `MaxPulse` values.
	ServoRangeDef    int     = 135
	ServoMinPulseDef float32 = 750.0
	ServoMaxPulseDef float32 = 2250.0
)

// Servo structure
type Servo struct {
	pca     *PCA9685
	channel int
	options *ServOptions
}

// ServOptions for servo
type ServOptions struct {
	AcRange  int // actuation range
	MinPulse float32
	MaxPulse float32
}

// ServoNew creates a new servo driver
func (pca *PCA9685) ServoNew(chn int, o *ServOptions) *Servo {
	s := &Servo{
		pca:     pca,
		channel: chn,
		options: &ServOptions{
			AcRange:  ServoRangeDef,
			MinPulse: ServoMinPulseDef,
			MaxPulse: ServoMaxPulseDef,
		},
	}
	if o != nil {
		s.options = o
	}
	return s
}

// Angle in degrees. Must be in the range `0` to `Range`.
func (s *Servo) Angle(a int) (err error) {
	if a < 0 || a > s.options.AcRange {
		return fmt.Errorf("Angle out of range")
	}
	return s.Fraction(float32(a) / float32(s.options.AcRange))
}

// Fraction as pulse width expressed between 0.0 `MinPulse` and 1.0 `MaxPulse`.
// For conventional servos, corresponds to the servo position as a fraction
// of the actuation range.
func (s *Servo) Fraction(f float32) (err error) {
	if f < 0.0 || f > 1.0 {
		return fmt.Errorf("Must be 0.0 to 1.0")
	}

	freq := s.pca.GetFreq()

	minDuty := s.options.MinPulse * freq / 1000000 * 0xFFFF
	maxDuty := s.options.MaxPulse * freq / 1000000 * 0xFFFF
	dutyRange := maxDuty - minDuty
	dutyCycle := (int(minDuty+f*dutyRange) + 1) >> 4

	return s.pca.SetChannel(int(s.channel), 0, dutyCycle)
}

// Reset channel
func (s *Servo) Reset() (err error) {
	return s.pca.SetChannel(s.channel, 0, 0)
}
