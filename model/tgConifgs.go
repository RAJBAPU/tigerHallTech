package models

import (
	"fmt"

	"github.com/astaxie/beego/orm"
)

type TgConfig struct {
	Id          int    `orm:"column(id);auto"`
	Key         string `orm:"column(key);size(100)" description:"Name of the config"`
	Value       string `orm:"column(value);size(256)" description:"Value of config"`
	Description string `orm:"column(description);size(256)"`
}

func (t *TgConfig) TableName() string {
	return "tg_config"
}

func init() {
	orm.RegisterModel(new(TgConfig))
}

type BeegoTgConfig struct{}

func (i *BeegoTgConfig) GetAllTgConfig() (v []TgConfig, err error) {
	o := orm.NewOrm()
	v = []TgConfig{}

	_, err = o.QueryTable(new(TgConfig)).All(&v)
	if err != nil {
		fmt.Println("error GetAllTgConfig: ", err)
		return nil, err
	}

	return
}

func (i *BeegoTgConfig) GetAllConfigs() map[string]string {

	var configs []TgConfig
	var err error

	configs, err = i.GetAllTgConfig()
	_configs := make((map[string]string))
	if err != nil {
		return _configs
	}

	for _, v := range configs {
		_configs[v.Key] = v.Value
	}
	return _configs
}
