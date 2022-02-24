package utils

import (
	"github.com/fogleman/gg"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	drivers "github.com/r4stl1n/micro-hal/code/pkg/drivers"
	base "github.com/r4stl1n/micro-hal/code/pkg/drivers/base"
)

type SSD1306Test struct {
}

func (cmd *SSD1306Test) Init() *SSD1306Test {
	*cmd = SSD1306Test{}

	return cmd
}

func (cmd *SSD1306Test) Command() *cobra.Command {
	return &cobra.Command{
		Use:                   "ssd1306",
		Aliases:               []string{"ssd"},
		Args:                  cobra.ExactArgs(1),
		ArgAliases:            []string{"i2c-address"},
		DisableFlagsInUseLine: true,
		Short:                 "test display",
		Run:                   cmd.Run,
	}
}

func (cmd *SSD1306Test) Run(_ *cobra.Command, args []string) {

	// We create a connection to the i2c interface on the raspberry pi
	logrus.Infof("Attempting to connect to the i2c address: %s", args[0])
	i2c, err := new(base.I2C).Init(drivers.DefaultSSD1306Address, args[0], base.DEFAULT_I2C_ADDRESS)

	if err != nil {
		logrus.Fatal(err)
	}

	// Next we create the needed driver to connect to the pca9685
	logrus.Info("Creating new connection to SSD1306Test")

	ssd, err := new(drivers.SSD1306).Init(i2c, nil)

	if err != nil {
		logrus.Fatal(err)
	}

	err = ssd.On()

	if err != nil {
		logrus.Fatal(err)
	}

	ssd.Clear()

	ctx := gg.NewContext(128, 64)

	ctx.SetRGB(0, 0, 0)
	ctx.Clear()
	ctx.SetRGB(1, 1, 1)
	ctx.DrawStringAnchored("HerpDerp", 0, 0, 0, 1)

	ssd.ShowImage(ctx.Image())

}
