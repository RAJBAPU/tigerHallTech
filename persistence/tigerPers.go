package persistence

import (
	models "simpl_pr/model"

	"github.com/astaxie/beego/orm"
)

type TgConfigPersistence interface {
	GetAllTgConfig() (v []models.TgConfig, err error)
	GetAllConfigs() map[string]string
}

var tgConfigPersistence = map[string]func() TgConfigPersistence{
	"mysql": mysqlTgConfig,
}

func NewTgConfigPersistence(runMode string) TgConfigPersistence {

	f, ok := tgConfigPersistence[runMode]
	if !ok {
		return nil
	}
	return f()
}

func mysqlTgConfig() TgConfigPersistence { return &models.BeegoTgConfig{} }

type TgTigerDetailsPersistence interface {
	GetTgTigerDetailsById(id int) (v *models.TgTigerDetails, err error)
	GetTgTigerDetailsByIdWithORM(id int, o orm.Ormer) (v *models.TgTigerDetails, err error)
	GetTgTigerDetails(name string, dob string) (v *models.TgTigerDetails, err error)
	GetAllAliveTigers() (v []*models.TgTigerDetails, err error)
	GetAllDeadTigers() (v []*models.TgTigerDetails, err error)
	GetAllTigers(offset int, limit int) (tigers []*models.TgTigerDetails, err error)
	UpdateTgTiger(data *models.TgTigerDetails, o orm.Ormer, updatedBy string, columns ...string) (err error)
	GetCountOfTigers() (total int, err error)
	AddTgTigerDetails(m *models.TgTigerDetails) (id int64, err error)
}

var tgTigerDetailsPersistence = map[string]func() TgTigerDetailsPersistence{
	"mysql": mysqlTgTigerDetails,
}

func NewTgTigerDetailsPersistence(runMode string) TgTigerDetailsPersistence {

	f, ok := tgTigerDetailsPersistence[runMode]
	if !ok {
		return nil
	}
	return f()
}

func mysqlTgTigerDetails() TgTigerDetailsPersistence { return &models.BeegoTgTigerDetails{} }

type TgTigerImagesPersistence interface {
	GetTgTigerImagesByI(id int) (v *models.TgTigerImages, err error)
	GetTgTigerImagesByIdWithORM(id int, o orm.Ormer) (v *models.TgTigerImages, err error)
	AddTgTigerImages(m *models.TgTigerImages) (id int64, err error)
	GetAllTigerSightings(tigerId int, offset int, limit int) (tigers []*models.TgTigerImages, err error)
	GetCountOfTigerSightings(tigerId int) (total int, err error)
}

var tgTigerImagesPersistence = map[string]func() TgTigerImagesPersistence{
	"mysql": mysqlTgTigerImages,
}

func NewTgTigerImagesPersistence(runMode string) TgTigerImagesPersistence {

	f, ok := tgTigerImagesPersistence[runMode]
	if !ok {
		return nil
	}
	return f()
}

func mysqlTgTigerImages() TgTigerImagesPersistence { return &models.BeegoTgTigerImages{} }
