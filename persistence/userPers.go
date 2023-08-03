package persistence

import (
	models "simpl_pr/model"

	"github.com/astaxie/beego/orm"
)

type TgUserPersistence interface {
	GetYpUserById(id int) (v *models.TgUser, err error)
	GetYpUserByIdWithORM(id int, o orm.Ormer) (v *models.TgUser, err error)
	UpdateTgUser(data *models.TgUser, o orm.Ormer, updatedBy string, columns ...string) (err error)
	GetUserByVerificationCode(verificationCode string) (v *models.TgUser, err error)
	GetUserByEmail(email string) (v *models.TgUser, err error)
	AddTgUser(m *models.TgUser) (id int64, err error)
}

var tgUserPersistence = map[string]func() TgUserPersistence{
	"mysql": mysqlTgUser,
}

func NewTgUsergPersistence(runMode string) TgUserPersistence {

	f, ok := tgUserPersistence[runMode]
	if !ok {
		return nil
	}
	return f()
}

func mysqlTgUser() TgUserPersistence { return &models.BeegoTgUser{} }
