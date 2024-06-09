package config

import (
	"fmt"
	"reflect"

	"github.com/spf13/viper"
)

type Config struct {
	Verbose bool
	Port    int
	Users   string
}

func DefaultConfig() *Config {
	return &Config{
		Port: 8080,
	}
}

func InitiateConfig(cfgFile string) (*Config, error) {
	cfg := DefaultConfig()

	viper.SetEnvPrefix("GB")
	viper.AutomaticEnv()
	setDefaultsFromStruct(cfg)

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
		if err := viper.ReadInConfig(); err == nil {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		}
	}

	// Unmarshal the config into the Config struct
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func setDefaultsFromStruct(s interface{}) {
	val := reflect.ValueOf(s).Elem()
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		defaultValue := field.Tag.Get("default")
		if defaultValue != "" {
			viper.SetDefault(field.Name, defaultValue)
		}
	}
}
