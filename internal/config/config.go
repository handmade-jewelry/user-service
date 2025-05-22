package config

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

const (
	GRPCPort             = "grpc.port"
	GRPCNetwork          = "grpc.network"
	HTTPPort             = "http.port"
	HTTPHost             = "http.host"
	HTTPGracefulTimeout  = "http.graceful_timeout"
	DBName               = "database.name"
	DBUser               = "database.user"
	DBPassword           = "database.password"
	DBHost               = "database.host"
	DBPort               = "database.port"
	SSLMode              = "database.ssl_mode"
	DBMaxCons            = "database.max_cons"
	DBMinCons            = "database.min_cons"
	DBMaxConLifetime     = "database.max_con_lifetime"
	VerificationTokenExp = "token.exp"
)

type Config struct {
	GRPCPort             string
	GRPCNetwork          string
	HTTPPort             string
	HTTPHost             string
	HTTPGracefulTimeout  time.Duration
	DBName               string
	DBUser               string
	DBPassword           string
	DbHost               string
	DbPort               uint16
	SSLMode              string
	DBMaxCons            int32
	DBMinCons            int32
	DBMaxConLifetime     time.Duration
	VerificationTokenExp time.Duration
}

func LoadConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("fatal error config file: %w", err)
	}

	return nil
}
