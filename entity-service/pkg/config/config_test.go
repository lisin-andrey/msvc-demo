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
	envkTimeout   = "DEMO_TIMEOUT"
	envkSvcPort   = "DEMO_SVCPORT"
	envkProvider  = "DEMO_PROVIDER"
	envkDbName    = "DEMO_POSTGRES_DBNAME"
	envkDbUser    = "DEMO_POSTGRES_USER"
	envkDbPass    = "DEMO_POSTGRES_PASS"
	envkDbHost    = "DEMO_POSTGRES_HOST"
	envkDbConnStr = "DEMO_POSTGRES_CONN_STRING"
	// envkDbName    = "DEMO.POSTGRES_DBNAME"
	// envkDbUser    = "DEMO_POSTGRES.USER"
	// envkDbPass    = "DEMO-POSTGRES_PASS"
	// envkDbHost    = "DEMO_POSTGRES-HOST"
	// envkDbConnStr = "DEMO-POSTGRES-CONN_STRING"
)

var envData = map[string]string{
	envkTimeout:   "100s",
	envkSvcPort:   "1234",
	envkProvider:  "some provider",
	envkDbName:    "some db",
	envkDbUser:    "some user",
	envkDbPass:    "some pwd",
	envkDbHost:    "some host",
	envkDbConnStr: "some coonecting string",
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

	const relConfigFilePath = "../entity-service.yaml"

	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	configFile := path.Join(path.Dir(dir), relConfigFilePath)

	port, err := strconv.Atoi(envData[envkSvcPort])
	duration, err := time.ParseDuration(envData[envkTimeout])

	etalonCfg := Config{
		ServicePort:  port,
		Timeout:      duration,
		ProviderType: envData[envkProvider],
		PgdDal: ConfigDalPostgres{
			DbName:           envData[envkDbName],
			UserName:         envData[envkDbUser],
			Password:         envData[envkDbPass],
			DbHost:           envData[envkDbHost],
			ConnectingString: envData[envkDbConnStr],
		},
	}

	cfg, err := ReadConfig(configFile, "DEMO")
	if assert.NoError(t, err, "ReadConfig failed") {
		assert.Equal(t, etalonCfg, *cfg)
	}
}
