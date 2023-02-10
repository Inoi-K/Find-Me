package config

import (
	"github.com/spf13/viper"
	"log"
)

var C config

type config struct {
	MainSphereCoefficient  float64
	OtherSphereCoefficient float64

	DatabaseURL string
}

func ReadConfig() {
	viper.SetConfigFile("../../pkg/config/config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("couldn't read config file: %v", err)
	}

	err = viper.Unmarshal(&C)
	if err != nil {
		log.Fatalf("couldn't unmarshal config file: %v", err)
	}
}
