package main

import (
	"log"

	"github.com/spf13/viper"
)

type Configuration struct {
	vpr          *viper.Viper
	DbConfig     DatabaseConfig `mapstructure:"database"`
	ServerConfig ServerConfig   `mapstructure:"server"`
}

type DatabaseConfig struct {
	FileName string `mapstructure:"fileName"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

func main() {
	vpr := viper.New()

	vpr.AddConfigPath(".")
	vpr.SetConfigName("viperConfig")
	err := vpr.ReadInConfig()

	if err != nil {
		log.Fatal(err)
	}

	log.Print(vpr.GetString("server.port"))

	conf := &Configuration{}
	err = vpr.Unmarshal(conf)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf(conf.ServerConfig.Port)
}
