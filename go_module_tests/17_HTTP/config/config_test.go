package configuration_test

import (
	"fmt"
	poker "learning/17_HTTP"
	configuration "learning/17_HTTP/config"
	"strconv"
	"testing"

	"github.com/spf13/viper"
)

const (
	fullFileName         string = "./temp_config.yaml"
	fileName             string = "temp_config"
	testConfigDbFileName string = "test.db.json"
	testConfigServerPort string = "5000"
)

var testConfig string = fmt.Sprintf(`
database:
   format: "json"
   fileName: %s
   name: "Postgres"
server:
   port: %s `, testConfigDbFileName, testConfigServerPort)

var defaultConfig map[string]interface{} = map[string]interface{}{
	"server.port":       8000,
	"database.fileName": "dbFileName",
	"database.name":     "MongoDb",
	"database.port":     1234,
}

type SpyViper struct {
	configFile         string
	configPath         string
	envPrefix          string
	automaticEnvCalled bool
	dbFileName         string
	serverPort         string
	dbName             string
}

func (s *SpyViper) Unmarshal(rawConf interface{}, opts ...viper.DecoderConfigOption) error {
	conf := rawConf.(configuration.Configuration)
	conf.SetDatabaseFileName(s.dbFileName)
	conf.SetServerPort(s.serverPort)

	return nil
}

func (s *SpyViper) GetString(key string) string {
	switch key {
	case "database.fileName":
		return s.dbFileName
	case "database.name":
		return s.dbName
	case "server.port":
		return s.serverPort
	default:
		return ""
	}
}

func (s *SpyViper) SetConfigName(configFile string) {
	s.configFile = configFile
}

func (s *SpyViper) SetDefault(key string, value interface{}) {
	switch key {
	case "database.fileName":
		s.dbFileName = value.(string)
	case "database.name":
		s.dbName = value.(string)
	case "server.port":
		s.serverPort = strconv.Itoa(value.(int))
	}
}

func (s *SpyViper) AddConfigPath(configPath string) {
	s.configPath = configPath
}

func (s *SpyViper) ReadInConfig() error {
	s.serverPort = testConfigServerPort
	s.dbFileName = testConfigDbFileName

	return nil
}

func (s *SpyViper) AutomaticEnv() {
	if s.envPrefix != "" {
		s.automaticEnvCalled = true
	}
}

func (s *SpyViper) SetEnvPrefix(prefix string) {
	s.envPrefix = prefix
}

func TestConfigurationRead(t *testing.T) {
	t.Run("Reads default config and config file when given nonempty string unit", func(t *testing.T) {
		vpr := &SpyViper{}
		conf := configuration.NewConfiguration(vpr)
		assertDefaultConfigLoadedUnmarshal(t, conf)
	})

	t.Run("Reads only default config when given empty string unit", func(t *testing.T) {
		vpr := &SpyViper{}
		conf := configuration.NewConfiguration(vpr)
		wantedFilePath := "."

		err := conf.Read(fileName, wantedFilePath, defaultConfig)

		poker.AssertNoError(t, err)
		assertDbName(t, vpr.dbName, defaultConfig["database.name"].(string))
		assertDbFileName(t, conf.GetDatabaseFileName(), testConfigDbFileName)

		assertAutomaticEnvCalled(t, vpr)
		assertConfigFileName(t, vpr, fileName, wantedFilePath)
		assertPort(t, vpr.serverPort, testConfigServerPort)
	})

	t.Run("Reads default config when given empty string", func(t *testing.T) {
		_, clean := poker.CreateTempFileOsOpenFile(t, testConfig, fullFileName)
		defer clean()
		conf := configuration.NewConfiguration(viper.New())
		assertDefaultConfigLoadedUnmarshal(t, conf)
	})

	t.Run("Reads config from file when give a non empty string", func(t *testing.T) {
		_, clean := poker.CreateTempFileOsOpenFile(t, testConfig, fullFileName)
		defer clean()

		conf := configuration.NewConfiguration(viper.New())
		err := conf.Read(fileName, ".", defaultConfig)

		poker.AssertNoError(t, err)
		assertDbName(t, conf.GetDatabaseFileName(), testConfigDbFileName)
	})
}

func TestReadV2(t *testing.T) {
	t.Run("Read loads default config when given an nil", func(t *testing.T) {
		vpr, err := configuration.ReadV2(nil, defaultConfig)
		assertDefaultConfigLoaded(t, err, vpr)
	})

	t.Run("Read overloads defaultConfig when given a config file", func(t *testing.T) {
		file, clean := poker.CreateTempFile(t, testConfig, fileName)
		defer clean()

		vpr, err := configuration.ReadV2(file, defaultConfig)

		poker.AssertNoError(t, err)
		assertDbName(t, vpr.GetString("database.name"), "Postgres")
		assertDbPort(t, vpr.GetInt("database.port"), 1234)
	})
}

func TestRead(t *testing.T) {
	t.Run("Read loads default config when given an empty string", func(t *testing.T) {
		vpr, err := configuration.Read("", "", defaultConfig)
		assertDefaultConfigLoaded(t, err, vpr)
	})

	t.Run("Read overloads defaultConfig when given a config file", func(t *testing.T) {
		_, clean := poker.CreateTempFileOsOpenFile(t, testConfig, fullFileName)
		defer clean()

		vpr, err := configuration.Read(fileName, ".", defaultConfig)

		poker.AssertNoError(t, err)
		assertDbName(t, vpr.GetString("database.name"), "Postgres")
		assertDbPort(t, vpr.GetInt("database.port"), 1234)
	})
}

func assertDefaultConfigLoadedUnmarshal(t *testing.T, conf configuration.Configuration) {
	t.Helper()
	err := conf.Read("", "", defaultConfig)

	poker.AssertNoError(t, err)
	want := defaultConfig["database.fileName"].(string)
	got := conf.GetDatabaseFileName()
	assertDbName(t, got, want)

	want = strconv.Itoa(defaultConfig["server.port"].(int))
	got = conf.GetServerPort()
	assertPort(t, got, want)
}

func assertDefaultConfigLoaded(t *testing.T, err error, read configuration.Reader) {
	t.Helper()
	poker.AssertNoError(t, err)

	want := defaultConfig["database.name"].(string)
	got := read.GetString("database.name")

	assertDbName(t, got, want)
}

func assertDbFileName(t *testing.T, got, want string) {
	t.Helper()

	asserStrings(t, got, want,
		fmt.Sprintf("Did not load configuration properly. DbFileName mismatch: Wanted %s but got %s", want, got))
}

func assertDbName(t *testing.T, got, want string) {
	t.Helper()

	asserStrings(t, got, want,
		fmt.Sprintf("Did not load configuration properly. DbName mismatch: Wanted %s but got %s", want, got))
}

func assertPort(t *testing.T, got, want string) {
	t.Helper()
	asserStrings(t, got, want,
		fmt.Sprintf("Did not load configuration properly. Port mismatch: Wanted %s but got %s", want, got))
}

func assertAutomaticEnvCalled(t *testing.T, vpr *SpyViper) {
	t.Helper()
	if !vpr.automaticEnvCalled {
		t.Fatalf("Invalid config file read: got %s but wanted %s", vpr.configFile, fileName)
	}
}

func assertConfigFileName(t *testing.T, vpr *SpyViper, wantedName, wantedFilePath string) {
	t.Helper()

	if vpr.configFile != fileName {
		t.Fatalf("Invalid config file read: got %s but wanted %s", vpr.configFile, fileName)
	}

	if vpr.configPath != wantedFilePath {
		t.Fatalf("Invalid config file read: got %s but wanted %s", vpr.configFile, fileName)
	}
}

func asserStrings(t *testing.T, got, want, errMsg string) {
	t.Helper()
	if got != want {
		t.Fatalf(errMsg)
	}
}

func assertDbPort(t *testing.T, got, want int) {
	t.Helper()

	if got != want {
		t.Fatalf("Did not load default configuration properly. Wanted %d but got %d", want, got)
	}
}
