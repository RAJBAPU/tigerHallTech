package service

import (
	"errors"
	"fmt"
	models "simpl_pr/model"
	util "simpl_pr/utils"
	"time"

	"github.com/spf13/cast"
)

// [...] SignIn User
func SignInUser(payload SignInInput) (token string, err error) {
	functionName := "service.VerifyEmail"
	user, err := models.GetUserByEmail(payload.Email)
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

	configs := models.GetAllConfigs()
	expiersIn := cast.ToInt(configs["TokenExpiresIn"])
	tokeSecret := configs["TokenSecret"]
	expiers := time.Duration(expiersIn) * time.Minute

	// Generate Token
	token, err = util.GenerateToken(expiers, user.Id, tokeSecret)
	if err != nil {
		return
	}

	return
}
