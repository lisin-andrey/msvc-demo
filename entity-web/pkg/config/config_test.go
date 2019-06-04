package config

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	envkTimeout       = "DEMO_TIMEOUT"
	envkAppPort       = "DEMO_APP_PORT"
	envkEntityRestURL = "DEMO_ENTITY_REST_URL"
)

var envData = map[string]string{
	envkTimeout:       "100s",
	envkAppPort:       "1234",
	envkEntityRestURL: "some rest url",
}

func prepareEnv() {
	for k, v := range envData {
		os.Setenv(k, v)
	}
}

func printEnv() {
	for k := range envData {
		v := os.Getenv(k)
		fmt.Println(k, v)
	}
}

func clearEnv() {
	for k := range envData {
		os.Setenv(k, "")
	}
}

func TestConfigWithEnvData(t *testing.T) {
	prepareEnv()
	printEnv()

	const relConfigFilePath = "../entity-web.yaml"

	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	configFile := path.Join(path.Dir(dir), relConfigFilePath)

	port, err := strconv.Atoi(envData[envkAppPort])
	duration, err := time.ParseDuration(envData[envkTimeout])

	etalonCfg := Config{
		AppPort:       port,
		Timeout:       duration,
		EntityRestURL: envData[envkEntityRestURL],
	}

	cfg, err := ReadConfig(configFile, "DEMO")
	if assert.NoError(t, err, "ReadConfig failed") {
		assert.Equal(t, etalonCfg, *cfg)
	}
}
