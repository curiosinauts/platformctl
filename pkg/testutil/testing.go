package testutil

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// InitConfig initializes config for testing
func InitConfig() {
	// Find home directory.
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	// Search config in home directory with name ".platformctl" (without extension).
	viper.AddConfigPath(home)
	viper.SetConfigType("yaml")
	viper.SetConfigName(".platformctl")

	viper.SetEnvPrefix("platformctl")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("error while reading config %s, err: %s", viper.GetViper().ConfigFileUsed(), err)
	}
}
