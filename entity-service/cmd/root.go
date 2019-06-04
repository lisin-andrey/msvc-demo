package cmd

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/lisin-andrey/msvc-demo/common/pkg/tools"
	"github.com/lisin-andrey/msvc-demo/entity-service/pkg/config"
	"github.com/pkg/errors"
)

const (
	envKeyConfigFile = "ENTITY_SVC_CONFIG_FILE"
	defConfigFile    = "./entity-service.yaml"
	envKeyEnvPrefix  = "ENTITY_SVC_ENV_PREFIX"
	defEnvPrefix     = "DEMO"
)

// App - application data
type App struct {
	Config *config.Config
	Router *mux.Router
	Server *http.Server

	RepEntityHandler IRepInstanceHandler
}

// setupSignalHandler listens for syscall signals and stops HTTP server.
func (a *App) setupSignalHandler() {
	c := make(chan os.Signal, 2)

	//  SIGINT   (val:2)  | Завершение                 | Сигнал прерывания (Ctrl-C) с терминала
	//  SIGKILL  (val:9)  | Завершение                 | Безусловное завершение
	//  SIGQUIT  (val:3)  | Завершение с дампом памяти | Сигнал «Quit» с терминала (Ctrl-\)
	//  SIGTERM  (val:15) | Завершение                 | Сигнал завершения (сигнал по умолчанию для утилиты kill)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT)
	go func() {
		sig := <-c
		tools.Infofln("Got signal to shutdown: %d", sig)
		if sig == syscall.SIGTERM || sig == syscall.SIGINT || sig == os.Interrupt {
			tools.Infoln("Safe Shutting down.")

			// Create a deadline to wait for.
			ctx, cancel := context.WithTimeout(context.Background(), a.Config.Timeout)
			defer cancel()

			err := a.Server.Shutdown(ctx)
			if err != nil {
				tools.Errorfln("Failed to shutdown server: %s", err.Error())
			}
			os.Exit(1)
		} else if sig == syscall.SIGKILL {
			tools.Infoln("Closing server.")
			err := a.Server.Close()
			if err != nil {
				tools.Errorfln("Failed to close server: %s", err.Error())
			}
			os.Exit(2)
		}
	}()
}

// Run - run rest service
func (a *App) Run() error {
	a.Server = &http.Server{
		Handler:      handlers.RecoveryHandler()(handlers.LoggingHandler(os.Stdout, a.Router)),
		Addr:         "0.0.0.0:" + strconv.Itoa(a.Config.ServicePort),
		WriteTimeout: a.Config.Timeout,
		ReadTimeout:  a.Config.Timeout,
	}
	a.setupSignalHandler()

	tools.Debugfln("Start listening on port %d", a.Config.ServicePort)
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

	app.RepEntityHandler = &RepInstanceHandler{Config: app.Config}

	fRegisterHandler := func(f RestHandleFunc) http.Handler {
		return &RestHandler{
			RepHandler:     app.RepEntityHandler,
			RestHandleFunc: f,
		}
	}
	app.Router = NewRouter(fRegisterHandler)

	return app, nil
}
