package configuration

import (
	"fmt"
	"io"
	"log"

	repo "learning/17_HTTP/config/viper"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

//Configuration is an abstract interface representing that that must be accessible from a config
type Configuration interface {
	SetServerPort(newPort string)
	SetDatabaseFileName(newFileName string)
	GetServerPort() string
	GetDatabaseFileName() string
	Read(configFileName, configFilePath string, defaultConfig repo.DefaultConfiguration) error
}

//ConfigurationImpl is a type that holds the data required for the application to run
type ConfigurationImpl struct {
	reader   repo.Reader
	Server   ServerConfiguration
	Database DatabaseConfiguration
}

//ServerConfiguration is holds the configuration needed by the server like port, etc
type ServerConfiguration struct {
	Port string
}

//DatabaseConfiguration stores the name of our file which we are using as a database
type DatabaseConfiguration struct {
	FileName string
}

//NewConfiguration creates a configuration with an empty viper
func NewConfiguration(vpr repo.Reader) Configuration {
	return &ConfigurationImpl{
		vpr,
		ServerConfiguration{},
		DatabaseConfiguration{},
	}
}

//GetDatabaseFileName returns the database file name
func (c *ConfigurationImpl) GetDatabaseFileName() string {
	return c.Database.FileName
}

//SetDatabaseFileName returns the database file name
func (c *ConfigurationImpl) SetDatabaseFileName(newFileName string) {
	c.Database.FileName = newFileName
}

//SetServerPort returns the database file name
func (c *ConfigurationImpl) SetServerPort(newPort string) {
	c.Server.Port = newPort
}

//GetServerPort returns the port the server is running on
func (c *ConfigurationImpl) GetServerPort() string {
	return c.Server.Port
}

//Read generates the server viper configuration. You can give a default configuration not loaded
//from a file by giving an empty string for a fileName or filePath.
func (c *ConfigurationImpl) Read(configFileName, configFilePath string,
	defaultConfig repo.DefaultConfiguration) error {

	log.Printf("Loading default configuration")
	c.reader.LoadDefaultConfiguration(defaultConfig)

	if configFileName != "" && configFilePath != "" {
		log.Printf("Loading configuration from file")
		c.reader.LoadFromFile(configFileName, configFilePath)
	}

	err := c.reader.Unmarshal(c)
	if err != nil {
		return errors.Wrap(err, "Reader failed to Unmarshal")
	}

	return nil
}

//Read generates the server viper configuration. You can give a default configuration not loaded
//from a file by giving an empty string for a fileName or filePath.
func Read(configFileName, configFilePath string, defaultConfig repo.DefaultConfiguration) (*viper.Viper, error) {
	vCfg := viper.New()

	loadDefaultConfiguration(vCfg, defaultConfig)

	if configFileName == "" || configFilePath == "" {
		log.Printf("Loading default configuration")
		return vCfg, nil
	}

	err := loadFromFile(vCfg, configFileName, configFilePath)

	if err != nil {
		return nil, err
	}

	return vCfg, nil
}

func loadDefaultConfiguration(vCfg *viper.Viper, defaultConfig repo.DefaultConfiguration) {
	for key, value := range defaultConfig {
		vCfg.SetDefault(key, value)
	}
}

func loadFromFile(vCfg *viper.Viper, configFileName, configFilePath string) error {
	vCfg.SetEnvPrefix("vpr")
	vCfg.SetConfigName(configFileName)
	vCfg.AddConfigPath(configFilePath)
	vCfg.AutomaticEnv()
	err := vCfg.ReadInConfig()

	if err != nil {
		errText := fmt.Sprintf("Error reading configuration file: %s with path %s, %v",
			configFileName, configFilePath, err)

		return errors.Wrap(err, errText)
	}

	return nil
}

//ReadV2 takes in an io.Reader and reads it
func ReadV2(reader io.Reader, defaultConfig repo.DefaultConfiguration) (*viper.Viper, error) {
	vCfg := viper.New()
	loadDefaultConfiguration(vCfg, defaultConfig)

	if reader == nil {
		log.Println("Loading default configuration")
		return vCfg, nil
	}

	return readFromFile(vCfg, reader)
}

func readFromFile(vCfg *viper.Viper, reader io.Reader) (*viper.Viper, error) {
	vCfg.SetEnvPrefix("vpr")
	vCfg.AutomaticEnv()
	//Needed or ReadConfig does not work
	vCfg.SetConfigType("yaml")
	err := vCfg.ReadConfig(reader)

	if err != nil {
		errText := fmt.Sprintf("Error reading configuration file: %v", err)
		return nil, errors.Wrap(err, errText)
	}

	return vCfg, nil
}
