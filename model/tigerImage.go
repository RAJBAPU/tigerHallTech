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
	SightedByUser      int       `orm:"column(sightedByUser);"`
}

func (t *TgTigerImages) TableName() string {
	return "tg_tiger_sighting"
}

func init() {
	orm.RegisterModel(new(TgTigerImages))
}

type BeegoTgTigerImages struct{}

func (tg *BeegoTgTigerImages) GetTgTigerImagesByI(id int) (v *TgTigerImages, err error) {
	o := orm.NewOrm()
	return tg.GetTgTigerImagesByIdWithORM(id, o)
}

func (tg *BeegoTgTigerImages) GetTgTigerImagesByIdWithORM(id int, o orm.Ormer) (v *TgTigerImages, err error) {
	v = &TgTigerImages{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}

	return nil, err
}

func (tg *BeegoTgTigerImages) AddTgTigerImages(m *TgTigerImages) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

func (tg *BeegoTgTigerImages) GetAllTigerSightings(tigerId int, offset int, limit int) (tigers []*TgTigerImages, err error) {
	o := orm.NewOrm()
	//v = []*TgTigerImages{}

	//	_, err = o.QueryTable(new(TgTigerDetails)).Filter("tigerId", tigerId).OrderBy("-lastSteen").All(&v)
	query := "SELECT * FROM tg_tiger_sighting where tigerId = ? ORDER BY lastSeen DESC LIMIT ?, ?;"
	_, err = o.Raw(query, tigerId, offset, limit).QueryRows(&tigers)
	if err != nil {
		return nil, err
	}

	return
}

func (tg *BeegoTgTigerImages) GetCountOfTigerSightings(tigerId int) (total int, err error) {
	o := orm.NewOrm()

	count, err := o.QueryTable(new(TgTigerImages)).Filter("tigerId", tigerId).Count()
	if err != nil {
		fmt.Println("Error in GetCountOfTigerSightings: ", err)
		return
	}

	return cast.ToInt(count), err
}
