package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type DbConfig struct {
	Host      string `yaml:"host" env:"HOST" env-required:"true"`
	Port      string `yaml:"port" env:"PORT" env-required:"true"`
	Name      string `yaml:"name" env:"NAME" env-required:"true"`
	Namespace string `yaml:"namespace" env:"NAMESPACE" env-default:"items"`
}

type ServerConfig struct {
	Host string `yaml:"host" env:"HOST" env-required:"true"`
	Port string `yaml:"port" env:"PORT" env-required:"true"`
}

type Config struct {
	Server ServerConfig `yaml:"server" env-prefix:"SERVER_"`
	DB     DbConfig     `yaml:"database" env-prefix:"DB_"`

	TTL time.Duration `yaml:"ttl" env:"TTL" env-default:"15m"`
}

func New() *Config {
	return &Config{}
}

// Load загружает переменные из локального yaml файла и далее перезаписывает полученные и записывает отсутствующие значения переменными окружения.
// Этим обеспечивается вариативность конфигурации
func (cfg *Config) Load() error {
	err := cleanenv.ReadConfig("config.yaml", cfg)
	if err == nil {
		return nil
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return err
	}

	return nil
}
