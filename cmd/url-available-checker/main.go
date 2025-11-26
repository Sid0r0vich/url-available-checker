package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	handlers "github.com/Sid0r0vich/url-available-checker/internal/handlers"
)

var (
	SERVER_ADDR = ":8082"
)

func main() {
	api := handlers.NewAPI()
	r := handlers.NewMux(api)

	s := &http.Server{
		Addr:    SERVER_ADDR,
		Handler: r,
	}

	stopServer := make(chan os.Signal, 1)
	signal.Notify(stopServer, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-stopServer
		fmt.Println("\nStopping server")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.Shutdown(ctx); err != nil {
			fmt.Printf("Server error: %v\n", err)
		}
	}()

	fmt.Printf("Server started on %s\n", SERVER_ADDR)
	if err := s.ListenAndServe(); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
