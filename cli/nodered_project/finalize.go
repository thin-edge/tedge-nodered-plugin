/*
Copyright © 2024 thin-edge.io <info@thin-edge.io>
*/
package nodered_project

import (
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/thin-edge/tedge-nodered-plugin/pkg/cli"
)

func NewFinalizeCommand(ctx cli.Cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "finalize",
		Short: "Finalize operation",
		RunE: func(cmd *cobra.Command, args []string) error {
			slog.Debug("Executing", "cmd", cmd.CalledAs(), "args", args)
			return nil
		},
	}
	return cmd
}
