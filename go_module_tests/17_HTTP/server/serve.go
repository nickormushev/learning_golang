package server

import (
	"context"
	poker "learning/17_HTTP"
	configuration "learning/17_HTTP/config"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

//Server is an abstraction of a http.server
type Server interface {
	Shutdown(ctx context.Context) error
	ListenAndServe() error
}

//Application holds the http server and the configuration of the app
type Application struct {
	server  Server
	config  configuration.Configuration
	dbClose func()
}

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

//CreateApplication creates a application with injected server, configuration and dbClose methods
func CreateApplication(conf configuration.Configuration, server Server, dbClose func()) *Application {
	return &Application{
		server,
		conf,
		dbClose,
	}
}

//CreateDefaultApplication creates an http server and calls listen and serve
func CreateDefaultApplication(conf configuration.Configuration) *Application {
	store, dbClose, err := poker.GenerateFileSystemPlayerStore(conf.GetDatabaseFileName())

	if err != nil {
		log.Fatalf("Could not generate FileSystem player store from file, %v", err)
	}

	game := poker.NewGame(store, poker.BlindAlerterFunc(poker.GenericAlerter))
	playerServer, err := poker.NewPlayerServer(store, game)

	if err != nil {
		log.Fatalf("Failed to create playerServer %v", err)
	}

	//Maybe inject the server so you can mock it or are theese tests pointless?
	server := &http.Server{
		Addr:    conf.GetServerPort(),
		Handler: playerServer,
	}

	return &Application{
		server,
		conf,
		dbClose,
	}
}

//Start executes LisendAndServer for the server and waits for SIGINT to initiate gracefull shutdown
func (a *Application) Start() {
	go func() {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error when starting to listen %v", err)
		}
	}()

	log.Printf("Server started at port %s", a.config.GetServerPort())

	ctx := GenerateContextWithSigint()

	<-ctx.Done()
	a.gracefullShutdown()
}

//GracefullShutdown gracefully shuts down the http server
func (a *Application) gracefullShutdown() {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	defer func() {
		a.dbClose()
		cancel()
	}()

	if err := a.server.Shutdown(timeoutCtx); err != nil {
		log.Fatalf("Failed to gracefully shutdown server %v", err)
	}

	log.Printf("Successful graceful shutdown of server")
}
