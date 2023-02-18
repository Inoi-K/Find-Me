package config

import (
	"github.com/spf13/viper"
	"log"
	"time"
)

var C *config

type config struct {
	SphereID    int64
	Token       string
	DatabaseURL string

	MainSphereCoefficient  float64
	OtherSphereCoefficient float64

	Separator string // ArgumentSeparator
	ParseMode string

	Timeout time.Duration

	ProfileHost string
	ProfilePort string
}

func ReadConfig() {
	if C != nil {
		return
	}

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
