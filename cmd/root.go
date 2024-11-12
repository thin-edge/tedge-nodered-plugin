/*
Copyright © 2024 thin-edge.io <info@thin-edge.io>
*/
package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thin-edge/tedge-nodered-plugin/cli/nodered_flow"
	"github.com/thin-edge/tedge-nodered-plugin/cli/nodered_project"
	"github.com/thin-edge/tedge-nodered-plugin/pkg/cli"
)

// Build data
var buildVersion string
var buildBranch string
var Name = "tedge-nodered-plugin"
var EnvPrefix = "NODERED"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     Name,
	Short:   "thin-edge.io nodered plugin",
	Version: fmt.Sprintf("%s (branch=%s)", buildVersion, buildBranch),
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return SetLogLevel()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	args := os.Args
	name := filepath.Base(args[0])

	for _, c := range rootCmd.Commands() {
		// TODO: Only include commands with given annotation
		// c.Annotations
		if name == c.Name() {
			slog.Debug("Calling multi-call binary.", "name", name, "args", args)
			rootCmd.SetArgs(append([]string{name}, args[1:]...))
			break
		}
	}

	err := rootCmd.Execute()
	if err != nil {
		switch err.(type) {
		case cli.SilentError:
			// Don't log error
			slog.Debug("Silent error.", "err", err)
		default:
			slog.Error("Command error", "err", err)
		}
		os.Exit(1)
	}
}

func SetLogLevel() error {
	value := strings.ToLower(viper.GetString("log_level"))
	slog.Debug("Setting log level.", "new", value)
	switch value {
	case "info":
		slog.SetLogLoggerLevel(slog.LevelInfo)
	case "debug":
		slog.SetLogLoggerLevel(slog.LevelDebug)
	case "warn":
		slog.SetLogLoggerLevel(slog.LevelWarn)
	case "error":
		slog.SetLogLoggerLevel(slog.LevelError)
	}
	return nil
}

func init() {
	cliConfig := cli.Cli{}
	cobra.OnInitialize(func() {
		cliConfig.OnInit(Name, EnvPrefix)
	})
	rootCmd.AddCommand(
		nodered_flow.NewCommand(cliConfig),
		nodered_project.NewCommand(cliConfig),
	)

	// Don't show usage on errors
	rootCmd.SilenceUsage = true

	rootCmd.PersistentFlags().String("log-level", "info", "Log level")
	rootCmd.PersistentFlags().StringVarP(&cliConfig.ConfigFile, "config", "c", "", "Configuration file")

	// viper.Bind
	_ = viper.BindPFlag("log_level", rootCmd.PersistentFlags().Lookup("log-level"))
}
