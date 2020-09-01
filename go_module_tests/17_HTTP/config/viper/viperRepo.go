package configurationrepo

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

//ViperReader is an interface that reads configuration from a file and maps it it to
//a Configuration structure
//type ViperReader interface {
//	Unmarshal(rawVal interface{}, opts ...viper.DecoderConfigOption) error
//	SetConfigName(configFile string)
//	AddConfigPath(configPath string)
//	ReadInConfig() error
//	GetString(key string) string
//	SetDefault(key string, value interface{})
//	AutomaticEnv()
//	SetEnvPrefix(prefix string)
//}
//	c.loadDefaultConfiguration(defaultConfig)
//
//	c.loadFromFile(configFileName, configFilePath)
//
//	c.reader.Unmarshal(c)

//DefaultConfiguration is a structure that respresents the config required to start the application
type DefaultConfiguration map[string]interface{}

//Reader represents a configuration reader that does all the reading from a file
type Reader interface {
	LoadDefaultConfiguration(defaultConfig DefaultConfiguration)
	LoadFromFile(configFileName, configFilePath string) error
	Unmarshal(rawValue interface{}) error
}

//ViperReader uses the viper library to implement the viper class
type ViperReader struct {
	vpr *viper.Viper
}

//NewViperReader is a default constructor for ViperReader
func NewViperReader() Reader {
	return &ViperReader{viper.New()}
}

//LoadDefaultConfiguration loads a DefaultConfiguration inside of a viper.Viper
func (v *ViperReader) LoadDefaultConfiguration(defaultConfig DefaultConfiguration) {
	for key, value := range defaultConfig {
		v.vpr.SetDefault(key, value)
	}
}

//Unmarshal uses the viper library to unmarshal the configuration loaded into it into a structure
func (v *ViperReader) Unmarshal(rawValue interface{}) error {
	err := v.vpr.Unmarshal(rawValue)

	if err != nil {
		return errors.Wrap(err, "viper failed to Unmarshal")
	}

	return nil
}

//LoadFromFile reads a file into the viperConfiguration
func (v *ViperReader) LoadFromFile(configFileName, configFilePath string) error {
	v.vpr.SetEnvPrefix("vpr")
	v.vpr.SetConfigName(configFileName)
	v.vpr.AddConfigPath(configFilePath)
	v.vpr.AutomaticEnv()
	err := v.vpr.ReadInConfig()

	if err != nil {
		errText := fmt.Sprintf("Error reading configuration file: %s with path %s, %v",
			configFileName, configFilePath, err)

		return errors.Wrap(err, errText)
	}

	return nil
}
