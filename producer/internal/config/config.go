package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env        string `yaml:"env"`
	HttpServer `yaml:"http_server"`
	Kafka      `yaml:"kafka"`
}

type HttpServer struct {
	Port        int           `yaml:"port"`
	Address     string        `yaml:"address"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}

type Kafka struct {
	BootstrapServers string `yaml:"bootstrap_servers"`
	MessageTopic     string `yaml:"message_topic"`
	FlushTimeout     int    `yaml:"flush_timeout"`

	DeliveryTimeout time.Duration `yaml:"delivery_timeout"`
	Acks            int           `yaml:"acks"`
}

const ConfigPathVar = "CONFIG_PATH"

func GetConfig() *Config {
	s := os.Getenv(ConfigPathVar)
	if s == "" {
		log.Fatal(ConfigPathVar, "environment variable not set")
	}

	if _, err := os.Stat(s); os.IsNotExist(err) {
		log.Fatal(ConfigPathVar, "config file does not exist")
	}

	var config Config
	if err := cleanenv.ReadConfig(s, &config); err != nil {
		log.Fatalf("cannot read config file: %s", err)
	}

	return &config
}
