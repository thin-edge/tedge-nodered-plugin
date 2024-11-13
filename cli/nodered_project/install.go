/*
Copyright © 2024 thin-edge.io <info@thin-edge.io>
*/
package nodered_project

import (
	"encoding/json"
	"io"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/thin-edge/tedge-nodered-plugin/pkg/cli"
	"github.com/thin-edge/tedge-nodered-plugin/pkg/nodered"
)

type InstallCommand struct {
	*cobra.Command

	CommandContext cli.Cli
	ModuleVersion  string
	File           string
}

type ProjectDescription struct {
	Repository string `json:"repo,omitempty"`
}

// installCmd represents the install command
func NewInstallCommand(ctx cli.Cli) *cobra.Command {
	command := &InstallCommand{
		CommandContext: ctx,
	}
	cmd := &cobra.Command{
		Use:   "install <MODULE_NAME>",
		Short: "Install a project",
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
	client := nodered.NewClientWithRetries(GetAPI())

	file, err := os.Open(c.File)
	if err != nil {
		return err
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	project := &ProjectDescription{}
	err = json.Unmarshal(b, &project)
	if err != nil {
		return err
	}

	projects, err := client.ProjectList()
	if err != nil {
		return err
	}
	projectName := args[0]
	exists := false
	for _, project := range projects.Projects {
		if project == projectName {
			exists = true
			break
		}
	}

	if exists {
		slog.Info("Updating existing project.", "name", projectName)
		if _, err := client.ProjectSetActive(projectName, true); err != nil {
			return err
		}
		if _, err := client.ProjectPull(projectName); err != nil {
			return err
		}
	}

	slog.Info("Cloning new project.", "name", projectName)
	if _, err := client.ProjectClone(projectName, project.Repository); err != nil {
		return err
	}
	slog.Info("Activating project.", "name", projectName)
	if _, err := client.ProjectSetActive(projectName, true); err != nil {
		return err
	}

	slog.Info("Installed module.", "name", projectName, "url", project.Repository)
	return nil
}
