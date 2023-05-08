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

	DatabasePoolMaxConnections int32
	Timeout                    time.Duration

	ProfileHost string
	ProfilePort string
	REngineHost string
	REnginePort string
	MatchHost   string
	MatchPort   string

	GatewayHost string
	GatewayPort string
	VerifyPath  string
	VerifyKey   string
	EmailFrom   string
	SMTPHost    string
	SMTPPort    int
	SMTPUser    string
	SMTPPass    string

	Separator     string // ArgumentSeparator
	ParseMode     string
	LikeButton    string
	DislikeButton string
	HighlightMark string

	TagsLimit              int
	Alpha                  float64
	MainSphereCoefficient  float64
	OtherSphereCoefficient float64
	University             string
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
