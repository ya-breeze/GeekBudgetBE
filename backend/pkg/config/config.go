//nolint:forbidigo // it's okay to use fmt in this file
package config

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	SeedUsers                     string `mapstructure:"seed_users" default:""`
	JWTSecret                     string `mapstructure:"jwt_secret" default:""`
	Verbose                       bool   `mapstructure:"verbose" default:"false"`
	Port                          int    `mapstructure:"port" default:"8080"`
	DBPath                        string `mapstructure:"dbpath" default:":memory:"`
	DisableImporters              bool   `mapstructure:"disableimporters" default:"false"`
	DisableCurrenciesRatesFetch   bool   `mapstructure:"disablecurrenciesratesfetch" default:"false"`
	CookieSecure                  bool   `mapstructure:"cookiesecure" default:"true"`
	MatcherConfirmationHistoryMax int    `mapstructure:"matcherconfirmationhistorymax" default:"10"`
	BankImporterFilesPath         string `mapstructure:"bankimporterfilespath" default:"bank-importer-files"`
	BackupPath                    string `mapstructure:"backup_path"           default:""`
	BackupInterval                string `mapstructure:"backup_interval"       default:"24h"`
	BackupMaxCount                int    `mapstructure:"backup_max_count"      default:"10"`
}

func InitiateConfig(cfgFile string) (*Config, error) {
	cfg := Config{}

	setDefaultsFromStruct(&cfg)
	viper.SetEnvPrefix("GB")
	viper.AutomaticEnv()

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
		if err := viper.ReadInConfig(); err == nil {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		}
	}

	// Unmarshal the config into the Config struct
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	if cfg.Verbose {
		fmt.Printf("Config: %+v\n", cfg)
	}

	return &cfg, nil
}

func setDefaultsFromStruct(s interface{}) {
	val := reflect.ValueOf(s).Elem()
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		defaultValue := field.Tag.Get("default")
		mapKey := field.Tag.Get("mapstructure")
		if mapKey == "" {
			mapKey = strings.ToLower(field.Name)
		}
		viper.SetDefault(mapKey, defaultValue)
		// AutomaticEnv doesn't reliably find keys with underscores (e.g. jwt_secret →
		// GB_JWT_SECRET), so bind each key explicitly.
		_ = viper.BindEnv(mapKey)
	}
}
