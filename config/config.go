package config

import (
	"github.com/spf13/viper"
	"strings"
	"sync"
)

type (
	Config struct {
		Server *Server
		Db     *Database
	}

	Server struct {
		Host string
		Port string
	}

	Database struct {
		Host     string
		Port     string
		User     string
		Password string
		Dbname   string
		SSLMode  string
		TimeZone string
	}
)

var (
	configInstance *Config
	once           sync.Once
)

func GetConfig() *Config {
	once.Do(func() {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("./")
		viper.AutomaticEnv()
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

		if err := viper.ReadInConfig(); err != nil {
			panic(err)
		}

		if err := viper.Unmarshal(&configInstance); err != nil {
			panic(err)
		}
	})
	return configInstance
}
