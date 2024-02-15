package config

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"strings"
)

type Postgres struct {
	Url string `mapstructure:"url"`
}

type SMTP struct {
	Host            string `mapstructure:"host"`
	Port            string `mapstructure:"port"`
	From            string `mapstructure:"from"`
	User            string `mapstructure:"user"`
	Password        string `mapstructure:"password"`
	ConfirmTemplate string `mapstructure:"confirm-template"`
}

type Config struct {
	Postgres Postgres `mapstructure:"postgres"`
	SMTP     SMTP     `mapstructure:"smtp"`
}

func ReadConfigFile() Config {
	viper.SetConfigName("meet")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/meet")
	viper.AddConfigPath("config")
	viper.SetEnvKeyReplacer(strings.NewReplacer(`.`, `_`))
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("reading config file (probably doesn't exists): %w", err))
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		panic(fmt.Errorf("unmarshaling config file: %w", err))
	}

	log.Debug().Interface("Configuration", config).Msg("")
	return config
}
