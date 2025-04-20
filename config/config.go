package config

import (
	"github.com/spf13/viper"
	"github.com/tesla59/blaze/types"
	"strings"
	"sync"
)

var (
	configInstance *types.Config
	once           sync.Once
)

func GetConfig() *types.Config {
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
