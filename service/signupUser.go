package service

import (
	"fmt"
	"strings"

	models "simpl_pr/model"
	"simpl_pr/persistence"
	util "simpl_pr/utils"

	"github.com/astaxie/beego/orm"
	"github.com/thanhpk/randstr"

	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB
}

func NewAuthController(DB *gorm.DB) AuthController {
	return AuthController{DB}
}

type User struct {
	Configs persistence.TgConfigPersistence
	User    persistence.TgUserPersistence
}

// [...] SignUp User
func (tg *User) SignUpUser(payload SignUpInput) (errCode int, err error) {
	functionName := "serive.SignUpUser"
	valid := util.IsValidEmail(payload.Email)
	if !valid {
		fmt.Println(functionName, "invalid Email or Password")
		errCode = 452
		return
	}

	tgUser, err := tg.User.GetUserByEmail(payload.Email)
	if (err != nil && err != orm.ErrNoRows) || tgUser != nil {
		fmt.Println(functionName, "email Already Exist", err)
		errCode = 453
		return errCode, nil
	}

	if payload.Password != payload.PasswordConfirm {
		fmt.Println(functionName, "invalid Email or Password")
		errCode = 452
		return
	}

	hashedPassword, err := util.HashPassword(payload.Password)
	if err != nil {
		fmt.Println(functionName, "Error in HashPassword", err)
		return
	}

	newUser := models.TgUser{
		Name:     payload.Name,
		Email:    strings.ToLower(payload.Email),
		Password: hashedPassword,
	}

	_, err = tg.User.AddTgUser(&newUser)
	if err != nil {
		fmt.Println(functionName, "Error in AddTgUser: ", err)
		return
	}

	// Generate Verification Code
	code := randstr.String(10)

	verification_code := util.Encode(code)

	// Update User in Database
	newUser.VerificationCode = verification_code
	err = tg.User.UpdateTgUser(&newUser, nil, "SignUpUser", "verificationCode")
	if err != nil {
		fmt.Println(" error UpdateTgUser: ", err)
		return
	}

	var firstName = newUser.Name

	if strings.Contains(firstName, " ") {
		firstName = strings.Split(firstName, " ")[1]
	}

	// Send Email
	emailData := util.EmailData{
		FirstName: firstName,
		Subject:   "Your account verification code",
		Text:      "Enter the below code to verify your account: " + code,
	}

	configs := tg.Configs.GetAllConfigs()
	emailFrom := configs["EmailFrom"]

	err = util.SendEmail(emailFrom, newUser.Email, &emailData)
	if err != nil {
		fmt.Println(functionName, "Could not Sent Mail", err)
		return
	}

	fmt.Println("We sent an email with a verification code to " + newUser.Email)
	return
}
