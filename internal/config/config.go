package config

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

const (
	GRPCPort            = "grpc.port"
	GRPCNetwork         = "grpc.network"
	HTTPPort            = "http.port"
	HTTPHost            = "http.host"
	HTTPGracefulTimeout = "http.graceful_timeout"
)

type Config struct {
	GRPCPort            string
	GRPCNetwork         string
	HTTPPort            string
	HTTPHost            string
	HTTPGracefulTimeout time.Duration
}

func LoadConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	err := viper.ReadInConfig()
	if err != nil {
		//todo panic?..
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	return nil
}
