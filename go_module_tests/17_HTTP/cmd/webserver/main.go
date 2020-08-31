package main

import (
	server "learning/17_HTTP/server"
)

const (
	configFileName string = "viperConfig"
	configFilePath string = "./config"
)

func main() {
	app := server.CreateDefaultApplication(configFileName, configFilePath)
	app.Start()
}
