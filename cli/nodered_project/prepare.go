/*
Copyright © 2024 thin-edge.io <info@thin-edge.io>
*/
package nodered_project

import (
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/thin-edge/tedge-nodered-plugin/pkg/cli"
	"github.com/thin-edge/tedge-nodered-plugin/pkg/nodered"
)

// prepareCmd represents the prepare command
func NewPrepareCommand(ctx cli.Cli) *cobra.Command {
	return &cobra.Command{
		Use:   "prepare",
		Short: "Prepare for install/removal",
		RunE: func(cmd *cobra.Command, args []string) error {
			slog.Debug("Executing", "cmd", cmd.CalledAs(), "args", args)

			// Check if the node-red project mode is enabled
			client := nodered.NewClientWithRetries(GetAPI())
			_, err := client.ProjectList()
			return err
		},
	}
}
