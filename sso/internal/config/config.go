package config

import (
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Env      string        `mapstructure:"env"`
	TokenTTL time.Duration `mapstructure:"token_ttl" env-required:"true"`
	GRPC     GrpcServer    `mapstructure:"grpc-server"`
	Db       Database      `mapstructure:"database"`
}

type GrpcServer struct {
	Port    int           `mapstructure:"port"`
	Timeout time.Duration `mapstructure:"timeout"`
}

type Database struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Port     int    `mapstructure:"port"`
	NameDB   string `mapstructure:"dbname"`
	Host     string `mapstructure:"host"`
}

func MustLoad() *Config {

	path := ".\\config"

	var c Config
	viper.SetConfigName("local")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)

	err := viper.ReadInConfig()
	if err != nil {
		panic("error read config: " + err.Error())
	}

	err = viper.Unmarshal(&c)
	if err != nil {
		panic("error unvarshal config: " + err.Error())
	}

	return &c
}

func MustLoadByPath(configPath string) *Config {

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("invalid config path: " + configPath)
	}

	var c Config
	viper.SetConfigName("local")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)

	err := viper.ReadInConfig()
	if err != nil {
		panic("error read config: " + err.Error())
	}

	err = viper.Unmarshal(&c)
	if err != nil {
		panic("error unvarshal config: " + err.Error())
	}

	return &c
}
