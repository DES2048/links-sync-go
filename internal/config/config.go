package config

import "github.com/ilyakaznacheev/cleanenv"

type Core struct {
	LogLevel string `yaml:"logLevel" env:"LOG_LEVEL" env-default:"info"`
}

type Server struct {
	Address string `yaml:"address" env:"SERVER_ADDRESS" env-default:":6500"`
	ApiKey  string `yaml:"apiKey" env:"API_KEY" env-required:"true"`
}

type Db struct {
	Url string `yaml:"url" env:"DB_URL" env-required:"true"`
}
type Config struct {
	Core   Core   `yaml:"core"`
	Server Server `yaml:"server"`
	Db     Db     `yaml:"db"`
}

func ReadConfig(path string) (*Config, error) {
	config := &Config{}

	err := cleanenv.ReadConfig(path, config)

	return config, err
}
