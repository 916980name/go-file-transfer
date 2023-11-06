package config

import (
	"file-transfer/pkg/log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	ENV_PREFIX         = "GO_FILE_TRANSFER"
	defaultConfigName  = "file-transfer.yaml"
	recommendedHomeDir = ".file-transfer"
)

func ReadConfig(cfgFile string) {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		viper.AddConfigPath(filepath.Join(home, recommendedHomeDir))
		viper.AddConfigPath(".")

		viper.SetConfigType("yaml")

		viper.SetConfigName(defaultConfigName)
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix(ENV_PREFIX)
	replacer := strings.NewReplacer(".", "_", "-", "_")
	viper.SetEnvKeyReplacer(replacer)

	if err := viper.ReadInConfig(); err != nil {
		log.Errorw("Error reading config file", "error", err)
	}

	if err := viper.ReadInConfig(); err == nil {
		log.Infow("Using config file:", "file", viper.ConfigFileUsed())
	}
}
