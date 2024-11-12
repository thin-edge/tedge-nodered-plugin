/*
Copyright © 2024 thin-edge.io <info@thin-edge.io>
*/
package nodered_flow

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/thin-edge/tedge-nodered-plugin/pkg/cli"
	"github.com/thin-edge/tedge-nodered-plugin/pkg/nodered"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type InstallCommand struct {
	*cobra.Command

	CommandContext cli.Cli
	ModuleVersion  string
	File           string
}

// installCmd represents the install command
func NewInstallCommand(ctx cli.Cli) *cobra.Command {
	command := &InstallCommand{
		CommandContext: ctx,
	}
	cmd := &cobra.Command{
		Use:   "install <MODULE_NAME>",
		Short: "Install a flow",
		Args:  cobra.ExactArgs(1),
		RunE:  command.RunE,
	}

	cmd.Flags().StringVar(&command.ModuleVersion, "module-version", "", "Software version to install")
	cmd.Flags().StringVar(&command.File, "file", "", "File")
	command.Command = cmd
	return cmd
}

func (c *InstallCommand) RunE(cmd *cobra.Command, args []string) error {
	slog.Debug("Executing", "cmd", cmd.CalledAs(), "args", args)

	moduleName := args[0]

	client := nodered.NewClient(GetAPI())

	file, err := os.Open(c.File)
	if err != nil {
		return err
	}
	defer file.Close()

	var flowsIn any
	b, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	// Edit the flow configuration and add the flow name and version to it
	node := gjson.ParseBytes(b)
	flowIndexes := make([]int64, 0)
	if node.IsArray() {
		node.ForEach(func(key, value gjson.Result) bool {
			if value.Get("type").String() == "tab" {
				flowIndexes = append(flowIndexes, key.Int())
			}
			return true
		})
	}

	var ob []byte
	for _, i := range flowIndexes {
		ob, err = sjson.SetBytes(b, fmt.Sprintf("%d.env.-1", i), nodered.FlowEnv{Name: "MODULE_NAME", Value: moduleName, Type: "str"})
		if err != nil {
			return err
		}
		b = ob
		ob, err = sjson.SetBytes(b, fmt.Sprintf("%d.env.-1", i), nodered.FlowEnv{Name: "MODULE_VERSION", Value: c.ModuleVersion, Type: "str"})
		if err != nil {
			return err
		}
		b = ob
	}

	err = json.Unmarshal(ob, &flowsIn)
	if err != nil {
		return err
	}

	resp, err := client.SetFlow("", flowsIn)
	if err != nil {
		return err
	}

	slog.Info("New revision.", "rev", resp)
	return nil
}
