/*
Copyright © 2024 thin-edge.io <info@thin-edge.io>
*/
package nodered_flow

import (
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/thin-edge/tedge-nodered-plugin/pkg/cli"
)

// prepareCmd represents the prepare command
func NewPrepareCommand(ctx cli.Cli) *cobra.Command {
	return &cobra.Command{
		Use:   "prepare",
		Short: "Prepare for install/removal",
		Run: func(cmd *cobra.Command, args []string) {
			slog.Debug("Executing", "cmd", cmd.CalledAs(), "args", args)
		},
	}
}
