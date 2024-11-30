package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	Core struct {
		LogLevel string `yaml:"logLevel" env:"LOG_LEVEL" env-default:"info"`
	} `yaml:"core"`
	Server struct {
		Address string `yaml:"address" env:"SERVER_ADDRESS" env-default:":6500"`
	} `yaml:"server"`
	Db struct {
		Url string `yaml:"url" env:"DB_URL,env-required"`
	} `yaml:"db"`
}

func ReadConfig(path string) (*Config, error) {
	config := &Config{}

	err := cleanenv.ReadConfig(path, config)

	return config, err
}