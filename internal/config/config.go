package config

import (
	"log"

	"github.com/spf13/viper"
)

type ApplicationConfig struct {
	SpotifyConfig spotifyConfig `mapstructure:"spotify"`
	JWTSigningKey string        `mapstructure:"jwt_signing_key"`
	DatabaseURL   string        `mapstructure:"database_url"`
}

func NewApplicationConfig() *ApplicationConfig {

	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln(err)
	}

	var c ApplicationConfig

	err = viper.Unmarshal(&c)
	if err != nil {
		panic(err)
	}

	return &c
}

type spotifyConfig struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	RedirectURI  string `mapstructure:"redirect_uri"`
}
