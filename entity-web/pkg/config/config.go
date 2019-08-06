package config

import (
	"strings"
	"time"

	"github.com/lisin-andrey/msvc-demo/common/pkg/tools"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// Config - configuration of the app
type Config struct {
	AppPort int           `mapstructure:"app-port"`
	Timeout time.Duration `mapstructure:"timeout"`

	EntityRestURL string `mapstructure:"entity-rest-url"`
}

// ReadConfig - read configuration
func ReadConfig(configFile string, envPrefix string) (*Config, error) {
	tools.Debugfln("Start read config: configFile=[%s] envPrefix=[%s]", configFile, envPrefix)
	v := viper.New()

	if configFile != "" {
		v.SetConfigFile(configFile)
	} else {
		v.AddConfigPath("./")
		v.SetConfigName("config-entity-web")
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
