package cmd

import (
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/lisin-andrey/msvc-demo/common/pkg/tools"
	"github.com/lisin-andrey/msvc-demo/common/pkg/web"
	"github.com/lisin-andrey/msvc-demo/entity-web/pkg/config"
	"github.com/pkg/errors"
)

const (
	envKeyConfigFile = "ENTITY_WEB_CONFIG_FILE"
	defConfigFile    = "./entity-web.yaml"
	envKeyEnvPrefix  = "ENTITY_WEB_ENV_PREFIX"
	defEnvPrefix     = "DEMO"
)

// App - application data
type App struct {
	Config *config.Config
	Router *mux.Router
	Server *http.Server
	// WebData
}

// setupSignalHandler listens for syscall signals and stops HTTP server.
func (a *App) setupSignalHandler() {
	web.SetupSignalHandler(a.Server, a.Config.Timeout)
}

// Run - run rest service
func (a *App) Run() error {
	a.Server = &http.Server{
		// Handler:      handlers.RecoveryHandler()(handlers.LoggingHandler(os.Stdout, a.Router)),
		Handler:      handlers.LoggingHandler(os.Stdout, a.Router),
		Addr:         "0.0.0.0:" + strconv.Itoa(a.Config.AppPort),
		WriteTimeout: a.Config.Timeout,
		ReadTimeout:  a.Config.Timeout,
	}
	a.setupSignalHandler()

	tools.Debugfln("Start listening on port %d", a.Config.AppPort)
	return a.Server.ListenAndServe()
}

// Prepare - prepare data to run rest service
func Prepare() (*App, error) {
	var err error

	app := &App{}

	// Read config
	configFile := os.Getenv(envKeyConfigFile)
	if configFile == "" {
		configFile = defConfigFile
	}
	envPrefix := os.Getenv(envKeyEnvPrefix)
	if envPrefix == "" {
		envPrefix = defEnvPrefix
	}

	app.Config, err = config.ReadConfig(configFile, envPrefix)
	if err != nil {
		return nil, errors.Wrapf(err, "can't read config")
	}

	// exec, err := model.NewEntityRestExecutor(app.Config.EntityRestURL)
	// webData := WebData{EntityRestExecutor: exec}
	webData, err := newWebData(app.Config.EntityRestURL)
	if err != nil {
		return nil, errors.Wrapf(err, "can't create WebData")
	}

	fRegisterHandler := func(f WebHandleFunc) http.Handler {
		return &WebHandler{
			webData:       *webData,
			WebHandleFunc: f,
		}
	}
	app.Router = NewRouter(fRegisterHandler)

	return app, nil
}
