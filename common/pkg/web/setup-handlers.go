package web

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lisin-andrey/msvc-demo/common/pkg/tools"
)

// SetupSignalHandler listens for syscall signals and stops HTTP server.
func SetupSignalHandler(s *http.Server, timeout time.Duration) {
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
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			err := s.Shutdown(ctx)
			if err != nil {
				tools.Errorfln("Failed to shutdown server: %s", err.Error())
			}
			os.Exit(1)
		} else if sig == syscall.SIGKILL {
			tools.Infoln("Closing server.")
			err := s.Close()
			if err != nil {
				tools.Errorfln("Failed to close server: %s", err.Error())
			}
			os.Exit(2)
		}
	}()
}
