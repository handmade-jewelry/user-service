package app

import (
	"fmt"
	"github.com/spf13/viper"
)

const (
	jwtTokenSecret     = "auth.jwt_token_secret"
	authTokenExpMin    = "auth.auth_token_expiry_minutes"
	refreshTokenExpMin = "auth.refresh_token_expiry_minutes"
)

func initConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	return nil
}
