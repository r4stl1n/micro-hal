package cmds

import (
	"github.com/r4stl1n/micro-hal/code/internal/hal-utilities/cmds/servos"
	"github.com/spf13/cobra"
)

type Servo struct {
}

func (cmd *Servo) Init() *Servo {
	*cmd = Servo{}

	return cmd
}

func (cmd *Servo) Command() *cobra.Command {
	command := &cobra.Command{
		Use:                   "servo",
		Aliases:               []string{"s"},
		DisableFlagsInUseLine: true,
		Short:                 "servo commands",
	}

	command.AddCommand(new(servos.Move).Init().Command())

	return command
}
