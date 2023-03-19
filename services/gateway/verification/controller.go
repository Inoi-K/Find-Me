package verification

import (
	"github.com/Inoi-K/Find-Me/pkg/config"
	"github.com/Inoi-K/Find-Me/services/gateway/session"
	"github.com/gin-gonic/gin"
	"github.com/thanhpk/randstr"
	"log"
	"net/http"
)

var (
	codeUser = make(map[string]int64)
)

func SendEmail(email string, userID int64) error {
	// Generate Verification Code
	code := randstr.String(20)
	verificationCode := Encode(code)
	codeUser[verificationCode] = userID

	err := send(&EmailData{
		Email:   email,
		URL:     "http://" + config.C.GatewayHost + ":" + config.C.GatewayPort + "/" + config.C.VerifyPath + "/" + code,
		Subject: "Your account verification code",
	})
	if err != nil {
		log.Printf("could not send email to %s: %v", email, err)
		return err
	}

	log.Printf("An email with a verification code sent to %s", email)
	return nil
}

func VerifyEmail(ctx *gin.Context) {
	code := ctx.Params.ByName(config.C.VerifyKey)
	verificationCodeControl := Encode(code)

	userID, ok := codeUser[verificationCodeControl]
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid verification code or user doesn't exists"})
		session.UserStateArg[userID] <- "not verified"
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Email verified successfully"})
	session.UserStateArg[userID] <- "verified"
}
