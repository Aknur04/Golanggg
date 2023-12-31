package main

import (
	"context" // New import
	"errors"  // New import
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func (app *application) serve() error {
	server := &http.Server{
		Addr:         ":" + strconv.Itoa(app.config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		s := <-quit

		app.logger.PrintInfo("caught signal", map[string]string{
			"signal": s.String(),
		})
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		shutdownError <- server.Shutdown(ctx)

		err := server.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		app.logger.PrintInfo("completing background tasks", map[string]string{
			"addr": server.Addr,
		})

		app.wg.Wait()
		shutdownError <- nil
	}()

	app.logger.PrintInfo("starting server", map[string]string{
		"addr": server.Addr,
		"env":  app.config.env,
	})
	err := server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	err = <-shutdownError
	if err != nil {
		return err
	}
	app.logger.PrintInfo("stopped server", map[string]string{
		"addr": server.Addr,
	})
	return nil

}
