package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env      string        `yaml:"env" env-default:"local"`
	TokenTTL time.Duration `yaml:"token_ttl" env-required:"true"`
	Database DatabaseConfig
	GRPC     GRPCConfig
	Redis    RedisConfig
}

type DatabaseConfig struct {
	Uri          string `yaml:"uri" env-required:"true"`
	DatabaseName string `yaml:"databaseName" env-required:"true"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port" env-default:"44044"`
	Timeout time.Duration `yaml:"timeout" env-default:"2h"`
}

type RedisConfig struct {
	Address string `yaml:"address" env-default:"127.0.0.1"`
}

func MustLoad() *Config {
	path := fetchConfigPath()

	if path == "" {
		panic("Config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("Config file does not exist: " + path)
	}

	var config Config

	if err := cleanenv.ReadConfig(path, &config); err != nil {
		panic("Failed to read config file: " + err.Error())
	}

	return &config
}

func fetchConfigPath() string {
	var result string

	flag.StringVar(&result, "config", "", "Path to the config file")
	flag.Parse()

	if result == "" {
		result = os.Getenv("CONFIG_PATH")
	}

	return result
}
