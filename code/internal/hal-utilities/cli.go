package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/chzyer/readline"
	"github.com/r4stl1n/micro-hal/code/internal/hal-utilities/cmds"
	"github.com/spf13/cobra"
)

const (
	UsageTemplate = `Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
 {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}
`
)

type CLI struct {
	rootCommand   *cobra.Command
	isInteractive bool
	Error         error
}

func (c *CLI) Init() *CLI {

	*c = *new(CLI)

	c.rootCommand = &cobra.Command{
		Use:                   "hal-utilities",
		Short:                 "command line tool providing micro-hal utilities",
		Version:               "1.0",
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		Run:                   c.run,
	}

	cobra.OnInitialize(func() {
		c.rootCommand.ResetFlags()
	})

	c.rootCommand.PersistentFlags().BoolVarP(&c.isInteractive, "interactive", "i", false, "interactive mode")

	c.rootCommand.AddCommand(new(cmds.Servo).Init().Command())
	c.rootCommand.AddCommand(new(cmds.Utils).Init().Command())
	return c
}

func (c *CLI) Run() *CLI {
	if c.Error != nil {
		return c
	}

	c.Error = c.rootCommand.Execute()
	if c.Error != nil {
		return c
	}

	c.run(nil, nil)

	return c
}

func (c *CLI) run(*cobra.Command, []string) {
	if !c.isInteractive {
		return
	}

	if c.Error != nil {
		return
	}

	c.rootCommand.Use = "\b"
	completer := readline.NewPrefixCompleter()
	for _, child := range c.rootCommand.Commands() {
		c.cobraToReadline(completer, child)
	}

	dir, homeError := os.UserHomeDir()
	if homeError != nil {
		c.Error = homeError
		return
	}

	shell, readlineError := readline.NewEx(&readline.Config{
		Prompt:            "\033[31m>>\033[0m ",
		HistoryFile:       filepath.Join(dir, ".micro_hal_utilities_history"),
		HistoryLimit:      1000,
		AutoComplete:      completer,
		InterruptPrompt:   "^C",
		EOFPrompt:         "quit",
		HistorySearchFold: true,
	})

	if readlineError != nil {
		c.Error = readlineError
		return
	}

	defer func() { _ = shell.Close() }()

	for {
		line, readError := shell.Readline()
		if readError != nil {
			c.Error = readlineError
			return
		}

		cmd, args, findError := c.rootCommand.Find(strings.Fields(line))
		if findError != nil {
			fmt.Printf("failed to find command: %s", findError.Error())
			continue
		}

		if cmd == nil {
			fmt.Printf("failed to find command")
			continue
		}

		if cmd == c.rootCommand {
			fmt.Printf("failed to find command")
			continue
		}

		flagsError := cmd.ParseFlags(args)
		if flagsError != nil {
			fmt.Printf("invalid flags: %s", flagsError.Error())
			continue
		}

		if cmd.Run == nil {
			_ = cmd.Help()
			continue
		}

		if !cmd.DisableFlagParsing {
			args = cmd.Flags().Args()
		}

		argsError := cmd.ValidateArgs(args)
		if argsError != nil {
			_ = cmd.Help()
			continue
		}

		cmd.Run(cmd, args)
	}
}

func (c *CLI) cobraToReadline(node readline.PrefixCompleterInterface, cmd *cobra.Command) {
	cmd.SetUsageTemplate(UsageTemplate)
	cmd.Use = c.usage(cmd)
	item := readline.PcItem(cmd.Use)
	node.SetChildren(append(node.GetChildren(), item))
	for _, child := range cmd.Commands() {
		c.cobraToReadline(item, child)
	}
}

func (c *CLI) usage(cmd *cobra.Command) string {
	words := make([]string, 0, len(cmd.ArgAliases)+1)
	words = append(words, cmd.Use)

	for _, name := range cmd.ArgAliases {
		words = append(words, "["+name+"]")
	}

	return strings.Join(words, " ")
}
