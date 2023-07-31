package models

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/spf13/cast"
)

type TgTigerImages struct {
	Id                 int       `orm:"column(id);auto"`
	TigerId            int       `orm:"column(tigerId);null"`
	LastSteenTimeStamp time.Time `orm:"column(lastSeen);null"`
	Longitude          float64   `orm:"column(longitude);null"`
	Latitude           float64   `orm:"column(latitude);null"`
	Image              string    `orm:"column(image);null"`
}

func (t *TgTigerImages) TableName() string {
	return "tg_tiger_sighting"
}

func init() {
	orm.RegisterModel(new(TgTigerImages))
}

func GetTgTigerImagesByI(id int) (v *TgTigerImages, err error) {
	o := orm.NewOrm()
	return GetTgTigerImagesByIdWithORM(id, o)
}

func GetTgTigerImagesByIdWithORM(id int, o orm.Ormer) (v *TgTigerImages, err error) {
	v = &TgTigerImages{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}

	return nil, err
}

func AddTgTigerImages(m *TgTigerImages) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

func GetAllTigerSightings(tigerId int, offset int, limit int) (tigers []*TgTigerImages, err error) {
	o := orm.NewOrm()
	//v = []*TgTigerImages{}

	//	_, err = o.QueryTable(new(TgTigerDetails)).Filter("tigerId", tigerId).OrderBy("-lastSteen").All(&v)
	query := "SELECT * FROM tg_tiger_images where tigerId = ? ORDER BY lastSeen DESC LIMIT ?, ?;"
	_, err = o.Raw(query, tigerId, offset, limit).QueryRows(&tigers)
	if err != nil {
		return nil, err
	}

	return
}

func GetCountOfTigerSightings(tigerId int) (total int, err error) {
	o := orm.NewOrm()

	count, err := o.QueryTable(new(TgTigerImages)).Filter("tigerId", tigerId).Count()
	if err != nil {
		fmt.Println("Error in GetCountOfTigerSightings: ", err)
		return
	}

	return cast.ToInt(count), err
}
