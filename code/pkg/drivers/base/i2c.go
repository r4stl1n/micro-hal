// Before usage you should load the i2c-dev kernel module
//
//      sudo modprobe i2c-dev
//
package i2c

import (
	"encoding/hex"
	"os"
	"syscall"

	"github.com/sirupsen/logrus"
)

// Options represents a connection to I2C-device.
type Options struct {
	addr uint8
	dev  string
	rc   *os.File
}

const DEFAULT_I2C_ADDRESS = 0x0703

// New opens a connection for I2C-device.
// SMBus (System Management Bus) protocol over I2C
// supported as well: you should preliminary specify
// register address to read from, either write register
// together with the data in case of write operations.
func New(addr uint8, dev string, i2cAddress uintptr) (*Options, error) {
	i2c := &Options{
		addr: addr,
		dev:  "/dev/i2c-0",
	}

	if dev != "" {
		i2c.dev = dev
	}

	f, err := os.OpenFile(dev, os.O_RDWR, 0600)
	if err != nil {
		return i2c, err
	}
	if err := ioctl(f.Fd(), i2cAddress, uintptr(addr)); err != nil {
		return i2c, err
	}

	i2c.rc = f
	return i2c, nil
}

// GetAddr return device occupied address in the bus.
func (o *Options) GetAddr() uint8 {
	return o.addr
}

// GetDev return full device name.
func (o *Options) GetDev() string {
	return o.dev
}

// READ SECTION

// ReadBytes read bytes from I2C-device.
// Number of bytes read correspond to buf parameter length.
func (o *Options) ReadBytes(buf []byte) (int, error) {
	n, err := o.rc.Read(buf)
	if err != nil {
		return n, err
	}
	logrus.Debugf("Read %d hex bytes: [%+v]", len(buf), hex.EncodeToString(buf))
	return n, nil
}

// ReadRegBytes read count of n byte's sequence from I2C-device
// starting from reg address.
func (o *Options) ReadRegBytes(reg byte, n int) ([]byte, int, error) {
	logrus.Debugf("Read %d bytes starting from reg 0x%0X...", n, reg)
	if _, err := o.WriteBytes([]byte{reg}); err != nil {
		return nil, 0, err
	}
	buf := make([]byte, n)
	c, err := o.ReadBytes(buf)
	if err != nil {
		return nil, 0, err
	}
	return buf, c, nil
}

// ReadRegU8 reads byte from I2C-device register specified in reg.
func (o *Options) ReadRegU8(reg byte) (byte, error) {
	if _, err := o.WriteBytes([]byte{reg}); err != nil {
		return 0, err
	}
	buf := make([]byte, 1)
	if _, err := o.ReadBytes(buf); err != nil {
		return 0, err
	}
	logrus.Debugf("Read U8 %d from reg 0x%0X", buf[0], reg)
	return buf[0], nil
}

// ReadRegU16BE reads unsigned big endian word (16 bits)
// from I2C-device starting from address specified in reg.
func (o *Options) ReadRegU16BE(reg byte) (uint16, error) {
	if _, err := o.WriteBytes([]byte{reg}); err != nil {
		return 0, err
	}
	buf := make([]byte, 2)
	if _, err := o.ReadBytes(buf); err != nil {
		return 0, err
	}
	w := uint16(buf[0])<<8 + uint16(buf[1])
	logrus.Debugf("Read U16 %d from reg 0x%0X", w, reg)
	return w, nil
}

// ReadRegU16LE reads unsigned little endian word (16 bits)
// from I2C-device starting from address specified in reg.
func (o *Options) ReadRegU16LE(reg byte) (uint16, error) {
	w, err := o.ReadRegU16BE(reg)
	if err != nil {
		return 0, err
	}
	// exchange bytes
	w = (w&0xFF)<<8 + w>>8
	return w, nil
}

// ReadRegS16BE reads signed big endian word (16 bits)
// from I2C-device starting from address specified in reg.
func (o *Options) ReadRegS16BE(reg byte) (int16, error) {
	if _, err := o.WriteBytes([]byte{reg}); err != nil {
		return 0, err
	}
	buf := make([]byte, 2)
	if _, err := o.ReadBytes(buf); err != nil {
		return 0, err
	}
	w := int16(buf[0])<<8 + int16(buf[1])
	logrus.Debugf("Read S16 %d from reg 0x%0X", w, reg)
	return w, nil
}

// ReadRegS16LE reads signed little endian word (16 bits)
// from I2C-device starting from address specified in reg.
func (o *Options) ReadRegS16LE(reg byte) (int16, error) {
	w, err := o.ReadRegS16BE(reg)
	if err != nil {
		return 0, err
	}
	// exchange bytes
	w = (w&0xFF)<<8 + w>>8
	return w, nil
}

// WRITE SECTION

// WriteBytes send bytes to the remote I2C-device. The interpretation of
// the message is implementation-dependent.
func (o *Options) WriteBytes(buf []byte) (int, error) {
	logrus.Debugf("Write %d hex bytes: [%+v]", len(buf), hex.EncodeToString(buf))
	return o.rc.Write(buf)
}

// WriteRegU8 writes byte to I2C-device register specified in reg.
func (o *Options) WriteRegU8(reg byte, value byte) error {
	buf := []byte{reg, value}
	if _, err := o.WriteBytes(buf); err != nil {
		return err
	}
	logrus.Debugf("Write U8 %d to reg 0x%0X", value, reg)
	return nil
}

// WriteRegU16BE writes unsigned big endian word (16 bits)
// value to I2C-device starting from address specified in reg.
func (o *Options) WriteRegU16BE(reg byte, value uint16) error {
	buf := []byte{reg, byte((value & 0xFF00) >> 8), byte(value & 0xFF)}
	if _, err := o.WriteBytes(buf); err != nil {
		return err
	}
	logrus.Debugf("Write U16 %d to reg 0x%0X", value, reg)
	return nil
}

// WriteRegU16LE writes unsigned little endian word (16 bits)
// value to I2C-device starting from address specified in reg.
func (o *Options) WriteRegU16LE(reg byte, value uint16) error {
	w := (value*0xFF00)>>8 + value<<8
	return o.WriteRegU16BE(reg, w)
}

// WriteRegS16BE writes signed big endian word (16 bits)
// value to I2C-device starting from address specified in reg.
func (o *Options) WriteRegS16BE(reg byte, value int16) error {
	buf := []byte{reg, byte((uint16(value) & 0xFF00) >> 8), byte(value & 0xFF)}
	if _, err := o.WriteBytes(buf); err != nil {
		return err
	}
	logrus.Debugf("Write S16 %d to reg 0x%0X", value, reg)
	return nil
}

// WriteRegS16LE writes signed little endian word (16 bits)
// value to I2C-device starting from address specified in reg.
func (o *Options) WriteRegS16LE(reg byte, value int16) error {
	w := int16((uint16(value)*0xFF00)>>8) + value<<8
	return o.WriteRegS16BE(reg, w)
}

// WriteRegU24BE writes unsigned big endian word (24 bits)
// value to I2C-device starting from address specified in reg.
func (v *Options) WriteRegU24BE(reg byte, value uint32) error {
	buf := []byte{reg, byte(value >> 16 & 0xFF), byte(value >> 8 & 0xFF), byte(value & 0xFF)}
	if _, err := v.WriteBytes(buf); err != nil {
		return err
	}
	logrus.Debugf("Write U24 %d to reg 0x%0X", value, reg)
	return nil
}

// WriteRegU32BE writes unsigned big endian word (32 bits)
// value to I2C-device starting from address specified in reg.
func (v *Options) WriteRegU32BE(reg byte, value uint32) error {
	buf := []byte{reg, byte(value >> 24 & 0xFF), byte(value >> 16 & 0xFF), byte(value >> 8 & 0xFF), byte(value & 0xFF)}
	if _, err := v.WriteBytes(buf); err != nil {
		return err
	}
	logrus.Debugf("Write U32 %d to reg 0x%0X", value, reg)
	return nil
}

// Close I2C-connection.
func (o *Options) Close() error {
	return o.rc.Close()
}

func ioctl(fd, cmd, arg uintptr) error {
	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, fd, cmd, arg, 0, 0, 0); err != 0 {
		return err
	}
	return nil
}
