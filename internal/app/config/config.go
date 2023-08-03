package config

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v9"
)

type Config struct {
	Environment             string `env:"ENVIRONMENT" envDefault:"development"`
	ServerAdd               string `env:"RUN_ADDRESS"`
	LogLevel                string `env:"LOG_LEVEL" envDefault:"info"`
	SecretKey               string `env:"SECRET_KEY" envDefault:"this-key-is-so-secret-nobody-can-guess-it-12345"`
	TokenExp                int    `env:"TOKEN_EXP" envDefault:"720"`
	AccrualSystemAddress    string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	AccrualOrderChannelSize int    `env:"ORDER_CHANNEL_SIZE" envDefault:"5"`
	DataBaseURI             string `env:"DATABASE_URI"`
	AccrualTickerPeriod     string `env:"ACCRUAL_TICKER_PERIOD"`
	AccrualRetryInterval    string `env:"ACCRUAL_RETRY_INTERVAL"`
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) InitConfig() {
	flag.StringVar(&c.ServerAdd, "a", "localhost:8080", "server address")
	flag.StringVar(&c.AccrualSystemAddress, "r", "http://localhost:8080", "accrual system address")
	flag.StringVar(&c.DataBaseURI, "d", "", "database dsn")
	flag.StringVar(&c.AccrualTickerPeriod, "t", "1", "accrual ticker period")
	flag.StringVar(&c.AccrualRetryInterval, "ri", "1", "accrual retry interval")
	//для локального запуска
	//flag.StringVar(&c.DataBaseURI, "d", "postgres://guest:guest@localhost:5432/loyalty", "database dsn")

	err := env.Parse(c)
	if err != nil {
		log.Fatal(err)
	}
}
