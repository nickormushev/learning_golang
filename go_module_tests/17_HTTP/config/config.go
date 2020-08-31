package configuration

import (
	"fmt"
	"io"
	"log"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

//DefaultConfiguration is a structure that respresents the config required to start the application
type DefaultConfiguration map[string]interface{}

//Reader is an interface that reads configuration from a file and maps it it to
//a Configuration structure
type Reader interface {
	Unmarshal(rawVal interface{}, opts ...viper.DecoderConfigOption) error
	SetConfigName(configFile string)
	AddConfigPath(configPath string)
	ReadInConfig() error
	GetString(key string) string
	SetDefault(key string, value interface{})
	AutomaticEnv()
	SetEnvPrefix(prefix string)
}

//Configuration is an abstract interface representing that that must be accessible from a config
type Configuration interface {
	SetServerPort(newPort string)
	SetDatabaseFileName(newFileName string)
	GetServerPort() string
	GetDatabaseFileName() string
	Read(configFileName, configFilePath string, defaultConfig DefaultConfiguration) error
}

//ConfigurationImpl is a type that holds the data required for the application to run
type ConfigurationImpl struct {
	reader   Reader
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
func NewConfiguration(vpr Reader) Configuration {
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
	defaultConfig DefaultConfiguration) error {

	log.Printf("Loading default configuration")
	c.loadDefaultConfiguration(defaultConfig)

	if configFileName != "" && configFilePath != "" {
		log.Printf("Loading configuration from file")
		c.loadFromFile(configFileName, configFilePath)
	}

	err := c.reader.Unmarshal(c)
	if err != nil {
		return errors.Wrap(err, "Could not unmarshal")
	}

	return nil
}

func (c *ConfigurationImpl) loadDefaultConfiguration(defaultConfig DefaultConfiguration) {
	for key, value := range defaultConfig {
		c.reader.SetDefault(key, value)
	}
}

func (c *ConfigurationImpl) loadFromFile(configFileName, configFilePath string) error {
	c.reader.SetEnvPrefix("vpr")
	c.reader.SetConfigName(configFileName)
	c.reader.AddConfigPath(configFilePath)
	c.reader.AutomaticEnv()
	err := c.reader.ReadInConfig()

	if err != nil {
		errText := fmt.Sprintf("Error reading configuration file: %s with path %s, %v",
			configFileName, configFilePath, err)

		return errors.Wrap(err, errText)
	}

	return nil
}

//Read generates the server viper configuration. You can give a default configuration not loaded
//from a file by giving an empty string for a fileName or filePath.
func Read(configFileName, configFilePath string, defaultConfig DefaultConfiguration) (*viper.Viper, error) {
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

func loadDefaultConfiguration(vCfg *viper.Viper, defaultConfig DefaultConfiguration) {
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
func ReadV2(reader io.Reader, defaultConfig DefaultConfiguration) (*viper.Viper, error) {
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
