package service

import (
	"errors"
	"fmt"
	models "simpl_pr/model"
	utils "simpl_pr/utils"
)

func VerifyEmail(verificationCode string) (err error) {
	functionName := "service.VerifyEmail"

	verification_code := utils.Encode(verificationCode)
	fmt.Println("verificationCode=", verification_code)

	updatedUser, err := models.GetUserByVerificationCode(verification_code)
	if err != nil {
		fmt.Println(functionName, "Error in GetUserByVerificationCode: ", err)
		return
	}
	if updatedUser.Verified {
		return errors.New("user already verified")
	}

	updatedUser.VerificationCode = ""
	updatedUser.Verified = true
	err = models.UpdateTgUser(updatedUser, nil, "VerifyEmail", "verificationCode", "verified")
	if err != nil {
		fmt.Println(functionName, "Error in UpdateTgUser: ", err)
		return
	}

	return

}
