package server_test

import (
	"context"
	poker "learning/17_HTTP"
	configuration "learning/17_HTTP/config"
	server "learning/17_HTTP/server"
	"syscall"
	"testing"
	"time"
)

type SpyServer struct {
	listenAndServeCalled bool
	shutdownCalled       bool
}

func (s *SpyServer) Shutdown(ctx context.Context) error {
	s.shutdownCalled = true
	return nil
}

func (s *SpyServer) ListenAndServe() error {
	s.listenAndServeCalled = true
	return nil
}

type SpyConfiguration struct {
	dbFileName string
	serverPort string
}

func (s *SpyConfiguration) SetDatabaseFileName(fileName string) {
	s.dbFileName = fileName
}

func (s *SpyConfiguration) SetServerPort(port string) {
	s.serverPort = port
}

func (s *SpyConfiguration) GetDatabaseFileName() string {
	return s.dbFileName
}

func (s *SpyConfiguration) GetServerPort() string {
	return s.serverPort
}

func (s *SpyConfiguration) Read(configFileName, configFilePath string,
	defaultConfig configuration.DefaultConfiguration) error {
	return nil
}

func TestCreateContext(t *testing.T) {
	t.Run("Context should be killed only after SIGINT", func(t *testing.T) {
		ctx := server.GenerateContextWithSigint()

		select {
		case <-ctx.Done():
			t.Errorf("Expected context to not be cancelled before SIGINT but it was")
		case <-time.After(5 * time.Millisecond):
		}

		syscall.Kill(syscall.Getpid(), syscall.SIGINT)

		select {
		case <-ctx.Done():
		case <-time.After(5 * time.Millisecond):
			t.Errorf("Expected context be cancelled afte SIGIN but was not")
		}
	})
}

func TestAppStart(t *testing.T) {
	conf := &SpyConfiguration{poker.TestDbFileName, poker.TestServerPort}
	srv := &SpyServer{}
	var closeDbCalled bool

	app := server.CreateApplication(conf, srv, func() { closeDbCalled = true })
	go func() {
		app.Start()
	}()

	poker.AssertFalse(t, srv.shutdownCalled)
	poker.AssertTrueWithRetry(t, &srv.listenAndServeCalled)

	syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	time.Sleep(10 * time.Millisecond)

	poker.AssertTrueWithRetry(t, &srv.shutdownCalled)
	poker.AssertTrueWithRetry(t, &closeDbCalled)
}
