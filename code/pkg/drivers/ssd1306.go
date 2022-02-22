package drivers

import (
	"fmt"
	"image"

	i2c "github.com/r4stl1n/micro-hal/code/pkg/drivers/base"
)

const DEFAULT_SSD1306_ADDRESS = 0x3c

const (
	// default values
	ssd1306Width        = 128
	ssd1306Height       = 64
	ssd1306ExternalVCC  = false
	ssd1306SetStartLine = 0x40
	// fundamental commands
	ssd1306SetComOutput0 = 0xC0
	ssd1306SetComOutput1 = 0xC1
	ssd1306SetComOutput2 = 0xC2
	ssd1306SetComOutput3 = 0xC3
	ssd1306SetComOutput4 = 0xC4
	ssd1306SetComOutput5 = 0xC5
	ssd1306SetComOutput6 = 0xC6
	ssd1306SetComOutput7 = 0xC7
	ssd1306SetComOutput8 = 0xC8
	ssd1306SetContrast   = 0x81
	// scrolling commands
	ssd1306ContinuousHScrollRight  = 0x26
	ssd1306ContinuousHScrollLeft   = 0x27
	ssd1306ContinuousVHScrollRight = 0x29
	ssd1306ContinuousVHScrollLeft  = 0x2A
	ssd1306StopScroll              = 0x2E
	ssd1306StartScroll             = 0x2F
	// adressing settings commands
	ssd1306SetMemoryAddressingMode = 0x20
	ssd1306ColumnAddr              = 0x21
	ssd1306PageAddr                = 0x22
	// hardware configuration commands
	ssd1306SetSegmentRemap0     = 0xA0
	ssd1306SetSegmentRemap127   = 0xA1
	ssd1306DisplayOnResumeToRAM = 0xA4
	ssd1306SetDisplayNormal     = 0xA6
	ssd1306SetDisplayInverse    = 0xA7
	ssd1306SetDisplayOff        = 0xAE
	ssd1306SetDisplayOn         = 0xAF
	// timing and driving scheme commands
	ssd1306SetDisplayClock      = 0xD5
	ssd1306SetPrechargePeriod   = 0xD9
	ssd1306SetVComDeselectLevel = 0xDB
	ssd1306SetMultiplexRatio    = 0xA8
	ssd1306SetComPins           = 0xDA
	ssd1306SetDisplayOffset     = 0xD3
	// charge pump command
	ssd1306ChargePumpSetting = 0x8D
)

type SSD1306DisplayBuffer struct {
	width    int
	height   int
	pageSize int
	buffer   []byte
}

// NewDisplayBuffer creates a new DisplayBuffer.
func (ssd1306DisplayBuffer *SSD1306DisplayBuffer) Init(width, height, pageSize int) *SSD1306DisplayBuffer {
	ssd1306DisplayBuffer = &SSD1306DisplayBuffer{
		width:    width,
		height:   height,
		pageSize: pageSize,
	}

	ssd1306DisplayBuffer.buffer = make([]byte, ssd1306DisplayBuffer.Size())

	return ssd1306DisplayBuffer
}

// Size returns the memory size of the display buffer.
func (ssd1306DisplayBuffer *SSD1306DisplayBuffer) Size() int {
	return (ssd1306DisplayBuffer.width * ssd1306DisplayBuffer.height) / ssd1306DisplayBuffer.pageSize
}

// Clear the contents of the display buffer.
func (ssd1306DisplayBuffer *SSD1306DisplayBuffer) Clear() {
	ssd1306DisplayBuffer.buffer = make([]byte, ssd1306DisplayBuffer.Size())
}

// SetPixel sets the x, y pixel with c color.
func (ssd1306DisplayBuffer *SSD1306DisplayBuffer) SetPixel(x, y, c int) {
	idx := x + (y/ssd1306DisplayBuffer.pageSize)*ssd1306DisplayBuffer.width
	bit := uint(y) % uint(ssd1306DisplayBuffer.pageSize)
	if c == 0 {
		ssd1306DisplayBuffer.buffer[idx] &= ^(1 << bit)
	} else {
		ssd1306DisplayBuffer.buffer[idx] |= (1 << bit)
	}
}

// Set sets the display buffer with the given buffer.
func (ssd1306DisplayBuffer *SSD1306DisplayBuffer) Set(buf []byte) {
	ssd1306DisplayBuffer.buffer = buf
}

type SSD1306Init struct {
	displayClock         byte
	multiplexRatio       byte
	displayOffset        byte
	startLine            byte
	chargePumpSetting    byte
	memoryAddressingMode byte
	comPins              byte
	contrast             byte
	prechargePeriod      byte
	vComDeselectLevel    byte
}

func (ssd1306Init *SSD1306Init) Init128X64() *SSD1306Init {

	ssd1306Init = &SSD1306Init{
		displayClock:         0x80,
		multiplexRatio:       0x3F,
		displayOffset:        0x00,
		startLine:            0x00,
		chargePumpSetting:    0x14, // 0x10 if external vcc is set
		memoryAddressingMode: 0x00,
		comPins:              0x12,
		contrast:             0xCF, // 0x9F if external vcc is set
		prechargePeriod:      0xF1, // 0x22 if external vcc is set
		vComDeselectLevel:    0x40,
	}

	return ssd1306Init
}

func (ssd1306Init *SSD1306Init) Init128X32() *SSD1306Init {

	ssd1306Init = &SSD1306Init{
		displayClock:         0x80,
		multiplexRatio:       0x1F,
		displayOffset:        0x00,
		startLine:            0x00,
		chargePumpSetting:    0x14, // 0x10 if external vcc is set
		memoryAddressingMode: 0x00,
		comPins:              0x02,
		contrast:             0x8F, // 0x9F if external vcc is set
		prechargePeriod:      0xF1, // 0x22 if external vcc is set
		vComDeselectLevel:    0x40,
	}

	return ssd1306Init
}

func (ssd1306Init *SSD1306Init) Init96X16() *SSD1306Init {

	ssd1306Init = &SSD1306Init{
		displayClock:         0x60,
		multiplexRatio:       0x0F,
		displayOffset:        0x00,
		startLine:            0x00,
		chargePumpSetting:    0x14, // 0x10 if external vcc is set
		memoryAddressingMode: 0x00,
		comPins:              0x02,
		contrast:             0x8F, // 0x9F if external vcc is set
		prechargePeriod:      0xF1, // 0x22 if external vcc is set
		vComDeselectLevel:    0x40,
	}

	return ssd1306Init
}

func (ssd1306Init *SSD1306Init) GetSequence(screenFlipped bool) []byte {

	segmentRemap := byte(ssd1306SetSegmentRemap127)
	comOutput := byte(ssd1306SetComOutput8)

	if screenFlipped {
		segmentRemap = ssd1306SetSegmentRemap0
		comOutput = ssd1306SetComOutput0
	}

	return []byte{
		ssd1306SetDisplayNormal,
		ssd1306SetDisplayOff,
		ssd1306SetDisplayClock, ssd1306Init.displayClock,
		ssd1306SetMultiplexRatio, ssd1306Init.multiplexRatio,
		ssd1306SetDisplayOffset, ssd1306Init.displayOffset,
		ssd1306SetStartLine | ssd1306Init.startLine,
		ssd1306ChargePumpSetting, ssd1306Init.chargePumpSetting,
		ssd1306SetMemoryAddressingMode, ssd1306Init.memoryAddressingMode,
		segmentRemap,
		comOutput,
		ssd1306SetComPins, ssd1306Init.comPins,
		ssd1306SetContrast, ssd1306Init.contrast,
		ssd1306SetPrechargePeriod, ssd1306Init.prechargePeriod,
		ssd1306SetVComDeselectLevel, ssd1306Init.vComDeselectLevel,
		ssd1306DisplayOnResumeToRAM,
		ssd1306SetDisplayNormal,
	}
}

// SSD1306 is a Driver for the PCA9685 16-channel 12-bit PWM/Servo controller
type SSD1306 struct {
	i2c           *i2c.I2C
	initSequence  *SSD1306Init
	options       *SSD1306Options
	displayBuffer *SSD1306DisplayBuffer
}

type SSD1306Options struct {
	Name        string
	Width       int
	Height      int
	PageSize    int
	ExternalVCC bool
}

// New creates the new PCA9685 driver with specified i2c interface and options
func (ssd1306 *SSD1306) Init(i2c *i2c.I2C, options *SSD1306Options) (*SSD1306, error) {

	adr := i2c.GetAddr()

	if i2c.GetAddr() == 0 {
		return nil, fmt.Errorf(`I2C device is not initiated`)
	}

	ssd1306 = &SSD1306{
		i2c: i2c,
		options: &SSD1306Options{
			Name:        "SSD1306" + fmt.Sprintf("-0x%x", adr),
			Width:       128,
			Height:      64,
			PageSize:    8,
			ExternalVCC: false,
		},
	}

	if options != nil {
		ssd1306.options = options
	}

	ssd1306.displayBuffer = new(SSD1306DisplayBuffer).Init(ssd1306.options.Width, ssd1306.options.Height, ssd1306.options.PageSize)

	initError := ssd1306.setupInitProcess()

	if initError != nil {
		return nil, initError
	}

	ssd1306.Off()

	err := ssd1306.commands(ssd1306.initSequence.GetSequence(false))

	if err != nil {
		return nil, err
	}

	err = ssd1306.commands([]byte{ssd1306ColumnAddr, 0, byte(ssd1306.options.Width) - 1})

	if err != nil {
		return nil, err
	}

	err = ssd1306.commands([]byte{ssd1306PageAddr, 0, (byte(ssd1306.options.Height / ssd1306.options.PageSize)) - 1})

	if err != nil {
		return nil, err
	}

	ssd1306.On()

	return ssd1306, nil
}

func (ssd1306 *SSD1306) setupInitProcess() error {

	// Construct the needed
	switch {
	case ssd1306.options.Width == 128 && ssd1306.options.Height == 64:
		ssd1306.initSequence = new(SSD1306Init).Init128X64()
	case ssd1306.options.Width == 128 && ssd1306.options.Height == 32:
		ssd1306.initSequence = new(SSD1306Init).Init128X32()
	case ssd1306.options.Width == 96 && ssd1306.options.Height == 16:
		ssd1306.initSequence = new(SSD1306Init).Init96X16()
	default:
		return fmt.Errorf("%dx%d resolution is unsupported, supported resolutions: 128x64, 128x32, 96x16", ssd1306.options.Width, ssd1306.options.Height)
	}

	// check for external vcc
	if ssd1306.options.ExternalVCC {
		ssd1306.initSequence.chargePumpSetting = 0x10
		ssd1306.initSequence.contrast = 0x9F
		ssd1306.initSequence.prechargePeriod = 0x22
	}

	return nil

}

func (ssd1306 *SSD1306) GetOptions() *SSD1306Options {
	return ssd1306.options
}

// On turns on the display.
func (ssd1306 *SSD1306) On() error {
	err := ssd1306.command(ssd1306SetDisplayOn)
	return err
}

// Off turns off the display.
func (ssd1306 *SSD1306) Off() error {
	err := ssd1306.command(ssd1306SetDisplayOff)
	return err
}

// Clear clears the display buffer.
func (ssd1306 *SSD1306) Clear() {
	ssd1306.displayBuffer.Clear()
}

// Set sets a pixel in the buffer.
func (ssd1306 *SSD1306) Set(x, y, c int) {
	ssd1306.displayBuffer.SetPixel(x, y, c)
}

// Reset clears display.
func (ssd1306 *SSD1306) Reset() error {
	if err := ssd1306.Off(); err != nil {
		return err
	}
	ssd1306.Clear()
	if err := ssd1306.On(); err != nil {
		return err
	}
	return nil
}

func (ssd1306 *SSD1306) SetContrast(contrast byte) error {

	err := ssd1306.commands([]byte{ssd1306SetContrast, contrast})

	return err
}

func (ssd1306 *SSD1306) ShowImage(img image.Image) (err error) {
	if img.Bounds().Dx() != ssd1306.options.Width || img.Bounds().Dy() != ssd1306.options.Height {
		return fmt.Errorf("image must match display width and height: %dx%d", ssd1306.options.Width, ssd1306.options.Height)
	}
	ssd1306.Clear()
	for y, w, h := 0, img.Bounds().Dx(), img.Bounds().Dy(); y < h; y++ {
		for x := 0; x < w; x++ {
			c := img.At(x, y)
			if r, g, b, _ := c.RGBA(); r > 0 || g > 0 || b > 0 {
				ssd1306.Set(x, y, 1)
			}
		}
	}
	return ssd1306.Display()
}

func (ssd1306 *SSD1306) Display() error {
	_, err := ssd1306.i2c.WriteBytes(append([]byte{0x40}, ssd1306.displayBuffer.buffer...))
	return err
}

func (ssd1306 *SSD1306) command(b byte) error {
	_, err := ssd1306.i2c.WriteBytes([]byte{0x80, b})
	return err
}

func (ssd1306 *SSD1306) commands(commands []byte) (err error) {
	var command []byte
	for _, d := range commands {
		command = append(command, []byte{0x80, d}...)
	}
	_, err = ssd1306.i2c.WriteBytes(command)
	return err
}
