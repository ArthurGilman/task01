package config

import (
	"os"
)

type Config struct {
	//Prefix string
	Kafka struct {
		Host string
		Port string
	}

	PostgreSQL struct {
		DSN   string
		Talbe string
	}

	HTTP struct {
		Host string
		Port string
	}

	Redis struct {
		Host     string
		Port     string
		Password string
		DB       string
	}
}

var conf *Config = &Config{}

func Init() {
	//conf.Prefix = os.Getenv("PREFIX")
	conf.Kafka.Host = os.Getenv("HOST")
	conf.Kafka.Port = os.Getenv("PORT")

	conf.PostgreSQL.DSN = os.Getenv("DSN")
	conf.PostgreSQL.Talbe = os.Getenv("TABLE")

	conf.HTTP.Host = os.Getenv("HTTPHOST")
	conf.HTTP.Port = os.Getenv("HTTPPORT")

	conf.Redis.Host = os.Getenv("REDISHOST")
	conf.Redis.Port = os.Getenv("REDISPORT")
	conf.Redis.Password = os.Getenv("REDISPASS")
	conf.Redis.DB = os.Getenv("REDISDB")
}

func Get() *Config {
	return conf
}
