package cli

import (
	"fmt"
	"log/slog"
	"os"
	"sort"
	"strings"

	"github.com/spf13/viper"
	"github.com/thin-edge/tedge-nodered-plugin/pkg/utils"
)

var LinuxConfigFilePath = "/etc/tedge-container-plugin/config.toml"

type SilentError error

type Cli struct {
	ConfigFile string
}

func (c *Cli) OnInit(name string, envPrefix string) {
	if c.ConfigFile != "" && utils.PathExists(c.ConfigFile) {
		// Use config file from the flag.
		viper.SetConfigFile(c.ConfigFile)
	} else {
		if home, err := os.UserHomeDir(); err == nil {
			// Add home directory.
			viper.AddConfigPath(home)
		}

		if utils.PathExists(LinuxConfigFilePath) {
			viper.SetConfigFile(LinuxConfigFilePath)
		} else {
			// Search config in home directory with name ".cobra" (without extension).
			viper.SetConfigType("json")
			viper.SetConfigType("toml")
			viper.SetConfigType("yaml")
			viper.SetConfigName(name)
		}
	}

	viper.SetEnvPrefix(envPrefix)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err == nil {
		slog.Info("Using config file", "path", viper.ConfigFileUsed())
	}
}

func (c *Cli) GetString(key string) string {
	return viper.GetString(key)
}

func (c *Cli) GetBool(key string) bool {
	return viper.GetBool(key)
}

func (c *Cli) PrintConfig() {
	keys := viper.AllKeys()
	sort.Strings(keys)
	for _, key := range keys {
		slog.Info("setting", "item", fmt.Sprintf("%s=%v", key, viper.Get(key)))
	}
}
