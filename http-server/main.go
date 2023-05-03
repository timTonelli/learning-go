package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	host = "localhost"
	port = 8000
)

func main() {
	app := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: createRouter(),
	}

	// Start listening on another go routine so the app isn't blocked
	go func() {
		fmt.Printf("Server listening on http://%s:%d\n", host, port)
		if err := app.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Error: %s", err)
		}
	}()

	// Create new context to listen for interrupt
	sigCtx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	// Block until interrupt is recieved
	<-sigCtx.Done()

	// Cancel context, resuming normal behavior of interrupt and notify user of shutdown
	stop()
	fmt.Println("\rGracefully shutting down...\nPress Ctrl+C again to force")

	// Start graceful shutdown, with a max timeout of 30 sec
	timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	if err := app.Shutdown(timeoutCtx); err != nil {
		log.Fatalf("Error: %s", err)
	}
}
