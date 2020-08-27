package server

import (
	"context"
	"fmt"
	poker "learning/17_HTTP"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

//GenerateContextWithSigint creates a context and waits for a sigint to cancell it
func GenerateContextWithSigint() context.Context {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		oscall := <-c
		log.Printf("System call %v", oscall)
		cancel()
	}()

	return ctx
}

//CreateServerAndServe creates an http server and calls listen and serve
func CreateServerAndServe(ctx context.Context, playerServer *poker.PlayerServer, port string) *http.Server {
	server := &http.Server{
		Addr:    port,
		Handler: playerServer,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error when starting to listen %v", err)
		}
	}()

	log.Printf("Server started at port %s", server.Addr)
	return server
}

//GracefullShutdown gracefully shuts down the http server
func GracefullShutdown(ctx context.Context, server *http.Server) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	defer cancel()

	if err = server.Shutdown(ctx); err != nil {
		return fmt.Errorf("Failed to gracefully shutdown server %v", err)
	}

	return nil
}
