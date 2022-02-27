package cmds

import (
	"github.com/r4stl1n/micro-hal/code/internal/hal-utilities/cmds/utils"
	"github.com/spf13/cobra"
)

type Utils struct {
}

func (cmd *Utils) Init() *Utils {
	*cmd = Utils{}

	return cmd
}

func (cmd *Utils) Command() *cobra.Command {
	command := &cobra.Command{
		Use:                   "utils",
		Aliases:               []string{"u"},
		DisableFlagsInUseLine: true,
		Short:                 "utility commands",
	}

	command.AddCommand(new(utils.LSM6DS3Test).Init().Command())
	command.AddCommand(new(utils.SSD1306Test).Init().Command())

	return command
}
