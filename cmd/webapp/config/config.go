package config

import (
	"fmt"
	"github.com/spf13/viper"
)

// Registry is for the configuration values.
var Registry *viper.Viper

// Set the configs
func Set() {
	viper.AddConfigPath(".") // optionally look for config in the working directory
	viper.SetConfigName("settings")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	Registry = viper.GetViper()
}
