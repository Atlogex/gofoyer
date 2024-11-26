package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type Config struct {
	Env          string     `yaml:"env" env-required:"local"`
	StoragePath  string     `yaml:"storage_path" env-required:"./storage/gofoyer.db"`
	GPRC         GRPCConfig `yaml:"gprc"`
	TokenTTL     string     `yaml:"token_ttl" env-default:"1h"`
	GRPCPort     int        `yaml:"grpc_port" env-default:"8045"`
	GRPCTimeout  string     `yaml:"grpc_timeout"`
	GRPCMaxConns int        `yaml:"grpc_max_conns"`
}

type GRPCConfig struct {
	Port     int    `yaml:"port"`
	Timeout  string `yaml:"timeout"`
	MaxConns int    `yaml:"max_conns"`
}

func MustLoad() *Config {
	path := fetchConfigPath()

	if path == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist - " + path)
	}

	var config Config
	if err := cleanenv.ReadConfig(path, &config); err != nil {
		panic(err)
	}

	return &config
}

// fetchConfigPath fetches path to config file. It checks command line flag
// `-config` and environment variable `CONFIG_PATH` and returns the first
// non-empty value. If both are empty, it returns an empty string.
func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
