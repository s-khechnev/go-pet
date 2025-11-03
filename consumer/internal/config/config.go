package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env        string `yaml:"env"`
	GrpcServer `yaml:"grpc"`
	Kafka      `yaml:"kafka"`
	DB
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

type DB struct {
	ConnString string
	Host       string `env:"DB_HOST"`
	Port       int    `env:"DB_PORT"`
	Username   string `env:"DB_USERNAME"`
	Password   string `env:"DB_PASSWORD"`
	Database   string `env:"DATABASE"`
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

	if err := godotenv.Load(); err != nil {
		log.Fatalf("failed to load .env: %v", err)
	}

	var config Config
	if err := cleanenv.ReadConfig(s, &config); err != nil {
		log.Fatalf("cannot read config file: %s", err)
	}

	config.DB.ConnString = fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		config.DB.Username,
		config.DB.Password,
		config.DB.Host,
		config.DB.Port,
		config.DB.Database,
	)

	return &config
}
