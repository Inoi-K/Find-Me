package server

import (
	"github.com/Inoi-K/Find-Me/pkg/config"
	"github.com/Inoi-K/Find-Me/services/gateway/verification"
	"github.com/gin-gonic/gin"
	"log"
)

func Start() {
	r := gin.Default()
	r.GET(config.C.VerifyPath+"/:"+config.C.VerifyKey, verification.VerifyEmail)
	err := r.Run(":" + config.C.GatewayPort)
	log.Printf("server listening at %v", config.C.GatewayPort)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
