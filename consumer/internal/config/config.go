package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env        string `yaml:"env"`
	GrpcServer `yaml:"http_server"`
	Kafka      `yaml:"kafka"`
}

type GrpcServer struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type Kafka struct {
	BootstrapServers string        `yaml:"bootstrap_servers"`
	MessageTopic     string        `yaml:"message_topic"`
	GroupId          string        `yaml:"group_id"`
	PollTimeout      time.Duration `yaml:"poll_timeout"`
	SessionTimeout   time.Duration `yaml:"session_timeout"`
	AutoOffsetReset  string        `yaml:"auto_offset_reset"`
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
