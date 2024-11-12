/*
Copyright © 2024 thin-edge.io <info@thin-edge.io>
*/
package nodered_project

import (
	"fmt"
	"log/slog"
	"sort"

	"github.com/spf13/cobra"
	"github.com/thin-edge/tedge-nodered-plugin/pkg/cli"
	"github.com/thin-edge/tedge-nodered-plugin/pkg/nodered"
)

// listCmd represents the list command
func NewListCommand(cliContext cli.Cli) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List nodered projects",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			slog.Debug("Executing", "cmd", cmd.CalledAs(), "args", args)

			client := nodered.NewClient(GetAPI())
			resp, err := client.ProjectList()
			if err != nil {
				// Don't fail the API is not ready yet
				slog.Warn("nodered api is not yet available.", "err", err)
				return nil
			}

			sort.Strings(resp.Projects)

			for _, name := range resp.Projects {
				// nodered only supports getting info for the active project
				if resp.Active == name {
					project, err := client.ProjectGet(name)
					if err != nil {
						return err
					}
					fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s\n", name, project.Version)
				} else {
					fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s\n", name, "inactive")
				}
			}
			return nil
		},
	}
}
