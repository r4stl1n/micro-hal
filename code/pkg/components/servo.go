package components

import (
	"fmt"

	pca9685 "github.com/r4stl1n/micro-hal/code/pkg/drivers"
	"github.com/sirupsen/logrus"
)

const (
	ServoRangeDef    int     = 180
	ServoMinPulseDef float32 = 500.0
	ServoMaxPulseDef float32 = 2500.0
)

// Servo structure
type Servo struct {
	pca     *pca9685.PCA9685
	channel int
	options *ServoOptions
}

// ServoOptions for servo
type ServoOptions struct {
	ActuationRange int // actuation range
	MinPulse       float32
	MaxPulse       float32
}

// ServoNew creates a new servo driver
func (servo *Servo) New(pca *pca9685.PCA9685, servoChannel int, serverOptions *ServoOptions) *Servo {
	servo = &Servo{
		pca:     pca,
		channel: servoChannel,
		options: serverOptions,
	}

	return servo
}

// Angle in degrees. Must be in the range `0` to `Range`.
func (servo *Servo) Angle(angle int) (err error) {
	if angle < 0 || angle > servo.options.ActuationRange {
		return fmt.Errorf("Angle out of range")
	}
	return servo.Fraction(float32(angle) / float32(servo.options.ActuationRange))
}

// Fraction as pulse width expressed between 0.0 `MinPulse` and 1.0 `MaxPulse`.
// For conventional servos, corresponds to the servo position as a fraction
// of the actuation range.
func (servo *Servo) Fraction(fraction float32) (err error) {
	if fraction < 0.0 || fraction > 1.0 {
		return fmt.Errorf("Must be 0.0 to 1.0")
	}

	freq := servo.pca.GetFreq()

	minDuty := servo.options.MinPulse * freq / 1000000 * 0xFFFF
	maxDuty := servo.options.MaxPulse * freq / 1000000 * 0xFFFF
	dutyRange := maxDuty - minDuty
	dutyCycle := (int(minDuty+fraction*dutyRange) + 1) >> 4

	logrus.Info(dutyCycle)
	return servo.pca.SetChannel(int(servo.channel), 0, dutyCycle)
}

// Reset channel
func (servo *Servo) Reset() (err error) {
	return servo.pca.SetChannel(servo.channel, 0, 0)
}
