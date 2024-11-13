/*
Copyright © 2024 thin-edge.io <info@thin-edge.io>
*/
package nodered_flow

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
		Short: "List nodered flows",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			slog.Debug("Executing", "cmd", cmd.CalledAs(), "args", args)

			client := nodered.NewClientWithoutRetries(GetAPI())
			resp, err := client.GetFlows()
			if err != nil {
				// Don't fail the API is not ready yet
				slog.Warn("nodered api is not yet available.", "err", err)
				return nil
			}
			flowIndexes := make(map[string]int)
			for i, flow := range resp {
				flowIndexes[flow.GetName()] = i
			}

			keys := make([]string, 0, len(flowIndexes))
			for key := range flowIndexes {
				keys = append(keys, key)
			}
			sort.Strings(keys)

			for _, name := range keys {
				flow := resp[flowIndexes[name]]
				fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s\n", name, flow.GetVersion())
			}

			return nil
		},
	}
}
