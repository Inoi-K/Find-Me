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

	Separator string // ArgumentSeparator
	ParseMode string

	DatabasePoolMaxConnections int32
	Timeout                    time.Duration

	ProfileHost string
	ProfilePort string
	REngineHost string
	REnginePort string
	MatchHost   string
	MatchPort   string

	TagsLimit              int
	MainSphereCoefficient  float64
	OtherSphereCoefficient float64
	Faculties              []string
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
