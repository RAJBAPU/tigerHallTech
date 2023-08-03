package middleware

import (
	models "simpl_pr/model"
	"simpl_pr/persistence"
)

type Customer interface {
	GetUser(id int) (user *models.TgUser, err error)
	GetConfigs() (configs map[string]string)
}

type CustomerSvc struct {
	Configs persistence.TgConfigPersistence
	User    persistence.TgUserPersistence
}

func CustomerService(
	Configs persistence.TgConfigPersistence,
	User persistence.TgUserPersistence,
) Customer {
	return &CustomerSvc{
		Configs: Configs,
		User:    User,
	}
}

func (tg *CustomerSvc) GetUser(id int) (user *models.TgUser, err error) {
	user, err = tg.User.GetYpUserById(id)
	if err != nil {
		return
	}
	return
}

func (tg *CustomerSvc) GetConfigs() (configs map[string]string) {
	configs = tg.Configs.GetAllConfigs()
	return
}
