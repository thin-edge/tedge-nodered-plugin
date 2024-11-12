/*
Copyright © 2024 thin-edge.io <info@thin-edge.io>
*/
package nodered_flow

import (
	"errors"
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
		Short: "Remove flows",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			slog.Debug("Executing", "cmd", cmd.CalledAs(), "args", args)

			client := nodered.NewClient(GetAPI())

			flows, err := client.GetFlows()
			if err != nil {
				return err
			}
			errs := make([]error, 0)
			for _, flow := range flows {
				slog.Info("Removing flow.", "id", flow.ID)
				_, err := client.DeleteFlow(flow.ID)
				errs = append(errs, err)
			}
			return errors.Join(errs...)
		},
	}
	cmd.Flags().StringVar(&command.ModuleVersion, "module-version", "", "Software version to remove")
	return cmd
}
