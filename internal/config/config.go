package config

import (
	"flag"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"path/filepath"
	"time"
)

// TODO: recheck without defaults
type Config struct {
	Env          string        `yaml:"env" env-default:"local"`
	StoragePath  string        `yaml:"storage_path" env-default:"./database/gofoyer.db"`
	GPRC         GRPCConfig    `yaml:"gprc"`
	TokenTTL     time.Duration `yaml:"token_ttl" env-default:"1h"`
	GRPCPort     int           `yaml:"grpc_port" env-default:"8144"`
	GRPCTimeout  time.Duration `yaml:"grpc_timeout" env-default:"1h"`
	GRPCMaxConns int           `yaml:"grpc_max_conns"`
}

type GRPCConfig struct {
	Port     int    `yaml:"port" env-default:"8144"`
	Timeout  string `yaml:"timeout"`
	MaxConns int    `yaml:"max_conns"`
}

func MustLoad() *Config {
	path := fetchConfigPath()

	if path == "" {
		panic("config path is empty")
	}

	return MustLoadBypath(path)
}

func MustLoadBypath(path string) *Config {
	//if _, err := os.Stat(path); os.IsNotExist(err) {
	//    panic("config file does not exist - " + path)
	//}

	dir, _ := filepath.Abs("")
	fullPath := filepath.Join(dir, path)
	fmt.Println("Configuration dir: ", dir)
	fmt.Println("Configuration path: ", fullPath)

	var config Config
	if err := cleanenv.ReadConfig(fullPath, &config); err != nil {
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
