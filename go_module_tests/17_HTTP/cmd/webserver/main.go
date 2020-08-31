package main

import (
	configuration "learning/17_HTTP/config"
	server "learning/17_HTTP/server"
	"log"

	"github.com/spf13/viper"
)

const (
	configFileName string = "viperConfig"
	configFilePath string = "./config"
)

func main() {
	appConfig := configuration.NewConfiguration(viper.New())
	err := appConfig.Read(configFileName, configFilePath, nil)

	if err != nil {
		log.Fatalf("Could not read startup configuration file %v", err)
	}

	app := server.CreateDefaultApplication(appConfig)
	app.Start()
}
