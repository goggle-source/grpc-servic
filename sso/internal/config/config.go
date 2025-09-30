package config

import (
	"flag"
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Env         string        `mapstructure:"env"`
	StoragePath string        `mapstructure:"storage_path" env-required:"true"`
	TokenTTL    time.Duration `mapstructure:"token_ttl" env-required:"true"`
	GRPC        GrpcServer    `mapstructure:"grpc-server"`
}

type GrpcServer struct {
	Port    int           `mapstructure:"port"`
	Timeout time.Duration `mapstructure:"timeout"`
}

func MustLoad() *Config {

	path := FindPathConfig()
	if path == "" {
		panic("error get config_path")
	}

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

func FindPathConfig() string {

	var res string

	flag.StringVar(&res, "config", "", "specify the path to the directory where configuration file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")

	}

	return res
}
