package service

import (
	"errors"
	"fmt"
	util "simpl_pr/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

// [...] SignIn User
func (tg *User) SignInUser(payload SignInInput, ctx *gin.Context) (token string, err error) {
	functionName := "service.VerifyEmail"
	user, err := tg.User.GetUserByEmail(payload.Email)
	if err != nil {
		fmt.Println(functionName, "Error in GetUserByEmail: ", err)
		return
	}

	if !user.Verified {
		fmt.Println(functionName, "User not verified", err)
		return "", errors.New("please verify your email")
	}
	if err = util.VerifyPassword(user.Password, payload.Password); err != nil {
		fmt.Println("invalid email or Password")
		return "", errors.New("invalid email or Password")
	}

	configs := tg.Configs.GetAllConfigs()
	expiersIn := cast.ToInt(configs["TokenExpiresIn"])
	tokeSecret := configs["TokenSecret"]
	expiers := time.Duration(expiersIn) * time.Minute

	// Generate Token
	token, err = util.GenerateToken(expiers, user.Id, tokeSecret)
	if err != nil {
		return
	}

	ctx.SetCookie("token", token, cast.ToInt(configs["TokenExpiresIn"]), "/", "localhost", false, true)
	return
}
