package config

import (
	"flag"
	"os"
	"sync"

	"github.com/keystop/yaDiploma.git/pkg/logger"
	"github.com/caarlos0/env/v6"
)

var Cfg *Config
var once sync.Once

type Config struct {
	servAddr       string
	dbConnString   string
	accrualAddress string
	appDir         string
}

func (c Config) ServAddr() string {
	return c.servAddr
}

func (c Config) DBConnString() string {
	return c.dbConnString
}

func (c Config) ProgramPath() string {
	return c.appDir
}

func (c Config) AccrualAddress() string {
	return c.accrualAddress
}

type EnvOptions struct {
	ServAddr       string `env:"RUN_ADDRESS"`
	DBConnString   string `env:"DATABASE_URI"`
	AccrualAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

//checkEnv for get options from env to default application options.
func (c *Config) checkEnv() {

	e := new(EnvOptions)
	e.DBConnString = "user=kseikseich password=112233 dbname=yap sslmode=disable"
	err := env.Parse(e)
	if err != nil {
		logger.Info("Ошибка чтения конфигурации из переменного окружения", err)
	}
	if len(e.ServAddr) != 0 {
		c.servAddr = e.ServAddr
	}
	if len(e.DBConnString) != 0 {
		c.dbConnString = e.DBConnString
	}
	if len(e.AccrualAddress) != 0 {
		c.accrualAddress = e.AccrualAddress
	}

}

func (c *Config) setDefault() {
	c.servAddr = "localhost:8080"
	c.dbConnString = "user=kseikseich password=112233 dbname=yap sslmode=disable"
	//c.dbConnString = "postgresql://postgres:postgres@postgres/praktikum?sslmode=disable"
	c.accrualAddress = "http://localhost:8082"
}

//setFlags for get options from console to default application options.
func (c *Config) setFlags() {
	flag.StringVar(&c.servAddr, "a", c.servAddr, "a server address string")
	flag.StringVar(&c.dbConnString, "d", c.dbConnString, "a db connection string")
	flag.StringVar(&c.accrualAddress, "r", c.accrualAddress, "a accrual system address")
	flag.Parse()
}

func createConfig() {
	Cfg = new(Config)

	appDir, err := os.Getwd()
	if err != nil {
		logger.Error(err)
	}

	Cfg.setDefault()
	Cfg.checkEnv()
	Cfg.setFlags()
	Cfg.appDir = appDir
	logger.Info("Создан config")
}

// NewDefOptions return obj like Options interfase.
func NewConfig() {
	once.Do(createConfig)
}
