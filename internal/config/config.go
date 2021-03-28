// Package config contains config DAO
package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// AppConfig is a container for application config
type AppConfig struct {
	API *API `mapstructure:"api"`
	Db  *DB  `mapstructure:"db"`
}

type API struct {
	Bind string `mapstructure:"bind"`
}

type DB struct {
	ConnString       string        `mapstructure:"conn_string"`
	MaxOpenConns     int           `mapstructure:"max_open_conns"`
	ConnMaxLifetime  time.Duration `mapstructure:"conn_max_lifetime"`
	MigrationDirPath string        `mapstructure:"migration_dir_path"`
	MigrationTable   string        `mapstructure:"migration_table"`
}

// GetAppConfig returns *Config
func GetAppConfig() (*AppConfig, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); nil != err {
		return nil, fmt.Errorf("unable to read config from file")
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	config := new(AppConfig)
	err := viper.Unmarshal(config)
	if nil != err {
		return nil, fmt.Errorf("unable to decode into struct, %v", err)
	}

	return config, nil
}
