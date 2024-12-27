package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/thesambayo/digillets-api/api/routes"
)

func (app *application) serve() error {
	// Declare an HTTP httpServer with some sensible timeout settings, which listens on the
	// port provided in the config struct and uses the serveMux we created above as the
	// handler.
	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.Port),
		Handler:      routes.Handlers(app.config, app.models, app.httpx),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Create a shutdownError channel.
	// We will use this to receive any errors returned by the graceful Shutdown() function.
	shutdownError := make(chan error)

	// Start a background goroutine to listen for SIGINT/SIGTERM signals.
	go func() {
		// Create a quitChannel channel which carries os.Signal values.
		quitChannel := make(chan os.Signal, 1)
		// Use signal.Notify() to listen for incoming SIGINT and SIGTERM signals and
		// relay them to the quit channel. Any other signals will not be caught by
		// signal.Notify() and will retain their default behavior.
		signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
		// Read the signal from the quit channel.
		//BLOCKING! This code will block until a signal is received.
		signalFromQuitChannel := <-quitChannel
		// the above 3 lines of code intercepts the signals: syscall.SIGINT, syscall.SIGTERM

		// Log a message to say that the signal has been caught.
		// call the String() method on the signal to get the signal name and include it in the log entry properties.
		app.logger.PrintInfo("shutting down server", map[string]string{
			"signal": signalFromQuitChannel.String(),
		})

		// Create a context with a 5-second timeout for graceful shutdown.
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Call Shutdown() on our server, passing in the context we just made.
		// Shutdown() will return nil if the graceful shutdown was successful, or an
		// error (which may happen because of a problem closing the listeners, or
		// because the shutdown didn't complete before the 5-second context deadline is
		// hit). We relay this return value to the shutdownError channel.
		// shutdownError <- server.Shutdown(ctx)

		// Call Shutdown() on the server like before👆🏼, but now we only send on the
		// shutdownError channel if it returns an error.
		err := httpServer.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		// Log a message to say that we're waiting for any background goroutines to
		// complete their tasks.
		app.logger.PrintInfo("completing background tasks", map[string]string{
			"addr": httpServer.Addr,
		})

		// Call Wait() to block until our WaitGroup counter is zero --- essentially
		// blocking until the background goroutines have finished. Then we return nil on
		// the shutdownError channel, to indicate that the shutdown completed without
		// any issues.
		app.wg.Wait()
		shutdownError <- nil
	}()

	// Start the HTTP server.
	app.logger.PrintInfo("starting server", map[string]string{
		"addr": httpServer.Addr,
		"env":  app.config.Env,
	})

	// Calling Shutdown() on our server will cause ListenAndServe() to immediately
	// return a http.ErrServerClosed error. So if we see this error, it is actually a
	// good thing and an indication that the graceful shutdown has started. So we check
	// specifically for this, only returning the error if it is NOT http.ErrServerClosed.
	err := httpServer.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	// Otherwise, we wait to receive the return value from Shutdown() on the
	// shutdownError channel. If return value is an error, we know that there was a
	// problem with the graceful shutdown, and we return the error.
	// i:e Wait for shutdown to complete.
	err = <-shutdownError
	if err != nil {
		return err
	}

	// At this point we know that the graceful shutdown completed successfully,
	// and we log a "stopped server" message.
	app.logger.PrintInfo("stopped server", map[string]string{
		"addr": httpServer.Addr,
	})

	return nil
}
