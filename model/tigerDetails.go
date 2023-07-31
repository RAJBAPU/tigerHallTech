package models

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/spf13/cast"
)

type TgTigerDetails struct {
	Id                 int       `orm:"column(id);auto"`
	Name               string    `orm:"column(name);size(64);null"`
	Dob                string    `orm:"column(dob);null"`
	LastSteenTimeStamp time.Time `orm:"column(lastSeen);null"`
	Longitude          float64   `orm:"column(longitude);null"`
	Latitude           float64   `orm:"column(latitude);null"`
	IsDead             bool      `orm:"column(isDead);null"`
}

func (t *TgTigerDetails) TableName() string {
	return "tg_tiger_details"
}

func init() {
	orm.RegisterModel(new(TgTigerDetails))
}

func GetTgTigerDetailsById(id int) (v *TgTigerDetails, err error) {
	o := orm.NewOrm()
	return GetTgTigerDetailsByIdWithORM(id, o)
}

func GetTgTigerDetailsByIdWithORM(id int, o orm.Ormer) (v *TgTigerDetails, err error) {
	v = &TgTigerDetails{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}

	return nil, err
}

func AddTgTigerDetails(m *TgTigerDetails) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

func GetTgTigerDetails(name string, dob string) (v *TgTigerDetails, err error) {

	o := orm.NewOrm()
	v = &TgTigerDetails{}
	err = o.QueryTable(new(TgTigerDetails)).Filter("name", name).Filter("dob", dob).Filter("isDead", 0).OrderBy("-id").One(v)
	if err != nil && err == orm.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		fmt.Println("error in GetTgTigerDetails ", err)
		return
	}

	return
}

func GetAllAliveTigers() (v []*TgTigerDetails, err error) {
	o := orm.NewOrm()
	v = []*TgTigerDetails{}

	_, err = o.QueryTable(new(TgTigerDetails)).Filter("isDead", 0).All(&v)
	if err != nil {
		fmt.Println("error GetAllAliveTigers", err)
		return nil, err
	}

	return
}

func GetAllDeadTigers() (v []*TgTigerDetails, err error) {
	o := orm.NewOrm()
	v = []*TgTigerDetails{}

	_, err = o.QueryTable(new(TgTigerDetails)).Filter("isDead", 1).All(&v)
	if err != nil {
		fmt.Println("error GetAllDeadTigers ", err)
		return nil, err
	}

	return
}

func GetAllTigers(offset int, limit int) (tigers []*TgTigerDetails, err error) {
	o := orm.NewOrm()

	query := "SELECT * FROM tg_tiger_details ORDER BY lastSeen DESC LIMIT ?, ?"

	_, err = o.Raw(query, offset, limit).QueryRows(&tigers)
	if err != nil {
		return nil, err
	}

	return
}

func UpdateTgTiger(data *TgTigerDetails, o orm.Ormer, updatedBy string, columns ...string) (err error) {
	if o == nil {
		o = orm.NewOrm()
	}
	_, err = o.Update(data, columns...)
	if err != nil {
		fmt.Println("Error UpdateTgTiger ", err)
	}
	return
}

func GetCountOfTigers() (total int, err error) {
	o := orm.NewOrm()

	count, err := o.QueryTable(new(TgTigerDetails)).Count()
	if err != nil {
		fmt.Println("Error in GetCountOfTigers: ", err)
		return
	}

	return cast.ToInt(count), err
}
