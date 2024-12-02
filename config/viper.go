package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
)

const (
	KeyEnv   = "ENV"
	EnvLocal = "local"
)

type viperConfig struct {
}

func NewViper() IViperConfig {
	viper.AddConfigPath("./config")
	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	env := os.Getenv(KeyEnv)
	if env != "" {
		viper.Set(KeyEnv, env)
	} else {
		viper.Set(KeyEnv, EnvLocal)
	}

	if viper.GetString(KeyEnv) == EnvLocal {
		if err := viper.ReadInConfig(); err != nil {
			panic(err)
		}
		fmt.Printf("Config loaded in local environment: %v\n", viper.AllSettings())
	}

	return &viperConfig{}
}

// Get returns the value of the key if it exists, otherwise it returns the default value
func (v viperConfig) Get(key string, defaultValue interface{}) interface{} {
	if viper.IsSet(key) {
		switch viper.Get(key).(type) {
		case int:
			return viper.GetInt(key)
		case string:
			return viper.GetString(key)
		case bool:
			return viper.GetBool(key)
		default:
			log.Printf("Unsupported type for key %s", key)
			return defaultValue
		}
	} else {
		log.Printf("Key %s not found", key)
		return defaultValue
	}
}

func (v viperConfig) GetInt(key string, defaultValue int) int {
	return v.Get(key, defaultValue).(int)
}

func (v viperConfig) GetString(key string, defaultValue string) string {
	return v.Get(key, defaultValue).(string)
}

func (v viperConfig) GetBool(key string, defaultValue bool) bool {
	return v.Get(key, defaultValue).(bool)
}
