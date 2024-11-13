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

type RemoveCommand struct {
	*cobra.Command

	ModuleVersion string
}

// removeCmd represents the remove command
func NewRemoveCommand(ctx cli.Cli) *cobra.Command {
	command := &RemoveCommand{}
	cmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			slog.Debug("Executing", "cmd", cmd.CalledAs(), "args", args)
			projectName := args[0]

			client := nodered.NewClientWithRetries(GetAPI())

			// Note: This will fail if the current project is active
			if err := client.ProjectDelete(projectName); err != nil {
				return err
			}
			slog.Info("Uninstalled project.", "name", projectName)
			return nil
		},
	}
	cmd.Flags().StringVar(&command.ModuleVersion, "module-version", "", "Software version to remove")
	return cmd
}
