package nodered_flow

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thin-edge/tedge-nodered-plugin/pkg/cli"
)

func GetAPI() string {
	v := viper.GetString("nodered.api")
	if v == "" {
		v = "http://127.0.0.1:1880"
	}
	return v
}

// NewCommand returns a cobra command for `nodered-flows` subcommands
func NewCommand(cmdCli cli.Cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nodered-flows",
		Short: "nodered-flows thin-edge.io software management plugin",
	}
	cmd.AddCommand(
		NewPrepareCommand(cmdCli),
		NewInstallCommand(cmdCli),
		NewRemoveCommand(cmdCli),
		NewUpdateListCommand(cmdCli),
		NewListCommand(cmdCli),
		NewFinalizeCommand(cmdCli),
	)
	return cmd
}
