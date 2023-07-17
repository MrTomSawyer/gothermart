package config

import (
	"log"

	"github.com/caarlos0/env/v9"
)

type Config struct {
	Environment string `env:"ENVIRONMENT" envDefault:"development"`
	Port        string `env:"PORT" envDefault:":8080"`
	LogLevel    string `env:"LOG_LEVEL" envDefault:"info"`
	SecretKey   string `env:"SECRET_KEY" envDefault:"secret-key"`
	TokenExp    int    `env:"TOKEN_EXP" envDefault:"3"`

	PostgresDSN    string `env:"PG_DSN" envDefault:"postgres://guest:guest@localhost:5432/loyalty"`
	UserTableName  string `env:"USER_TABLE_NAME" envDefault:"users"`
	OrderTableName string `env:"ORDER_TABLE_NAME" envDefault:"orders"`
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) InitConfig() {
	err := env.Parse(c)
	if err != nil {
		log.Fatal(err)
	}
}
