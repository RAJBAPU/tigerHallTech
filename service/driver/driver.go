package driver

import (
	models "simpl_pr/model"
	"simpl_pr/persistence"
	svc "simpl_pr/service"

	"github.com/gin-gonic/gin"
)

type Tiger interface {
	PostTigerDetails(payload svc.TigerDetails) error
	GetAllTigers(page int, pageSize int) (response *svc.TigerDetailsWithPaginationResponse, err error)
	PostSightingDetails(request svc.TigerDetails, user *models.TgUser) (errorCode int, err error)
	GetSightingDetails(tigerId int, page int, pageSize int) (response *svc.TigerSightingResponse, err error)
}

type User interface {
	SignUpUser(newUser svc.SignUpInput) (errCode int, err error)
	SignInUser(payload svc.SignInInput, ctx *gin.Context) (token string, err error)
	GetUserDetails(currentUser *models.TgUser) (userResponse *svc.UserResponse, err error)
	VerifyEmail(verificationCode string) (err error)
}

func NewTiger(
	Configs persistence.TgConfigPersistence,
	User persistence.TgUserPersistence,
	TigerDetails persistence.TgTigerDetailsPersistence,
	TigerImages persistence.TgTigerImagesPersistence,

) Tiger {
	return &svc.Tiger{
		Configs:      Configs,
		User:         User,
		TigerDetails: TigerDetails,
		TigerImages:  TigerImages,
	}
}

func NewUser(
	Configs persistence.TgConfigPersistence,
	User persistence.TgUserPersistence,

) User {
	return &svc.User{
		Configs: Configs,
		User:    User,
	}
}
