package config

import (
	"strings"
	"time"

	"github.com/lisin-andrey/msvc-demo/common/pkg/tools"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// Config - configuration of the service
type Config struct {
	ServicePort  int               `mapstructure:"svcport"`
	Timeout      time.Duration     `mapstructure:"timeout"`
	ProviderType string            `mapstructure:"provider"`
	PgdDal       ConfigDalPostgres `mapstructure:"postgres"`
}

// ConfigDalPostgres - configuration for Postgres dal
type ConfigDalPostgres struct {
	ConnectingString string `mapstructure:"conn-string"`
	DbName           string `mapstructure:"dbname"`
	UserName         string `mapstructure:"user"`
	Password         string `mapstructure:"pass"`
	DbHost           string `mapstructure:"host"`
}

// ReadConfig - read configuration
func ReadConfig(configFile string, envPrefix string) (*Config, error) {
	v := viper.New()

	if configFile != "" {
		v.SetConfigFile(configFile)
	} else {
		v.AddConfigPath("./")
		v.SetConfigName("config-entity-svc")
	}

	tools.Debugfln("ConfigFileUsed: %s", v.ConfigFileUsed())
	if err := v.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "cannot reading config")
	}
	tools.Debugln("Config reading")

	if envPrefix != "" {
		tools.Debugfln("ConfigEnvPrefix: %s", envPrefix)
		v.SetEnvPrefix(envPrefix)
	}
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	v.AutomaticEnv()

	config := &Config{}
	err := v.Unmarshal(config)
	if err != nil {
		return nil, errors.Wrap(err, "cannot unmarshal config")
	}

	tools.Debugfln("Config: %+v", config)

	return config, nil
}
