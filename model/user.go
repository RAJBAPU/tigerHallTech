package models

import (
	"fmt"

	"github.com/astaxie/beego/orm"
)

type TgUser struct {
	Id               int    `orm:"column(id);auto"`
	Password         string `orm:"column(password);size(512);null"`
	Email            string `orm:"column(email);size(64);null"`
	Name             string `orm:"column(name);size(64);null"`
	Verified         bool   `orm:"column(verified)"`
	VerificationCode string `orm:"column(verificationCode);size(64);null"`
}

func (t *TgUser) TableName() string {
	return "tg_user"
}

func init() {
	orm.RegisterModel(new(TgUser))
}

// GetYpUserById retrieves YpUser by Id. Returns error if
// Id doesn't exist
func GetYpUserById(id int) (v *TgUser, err error) {
	o := orm.NewOrm()
	return GetYpUserByIdWithORM(id, o)
}

func GetYpUserByIdWithORM(id int, o orm.Ormer) (v *TgUser, err error) {
	v = &TgUser{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}

	return nil, err
}

func AddTgUser(m *TgUser) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

func UpdateTgUser(data *TgUser, o orm.Ormer, updatedBy string, columns ...string) (err error) {
	if o == nil {
		o = orm.NewOrm()
	}
	err = updateRowByColumns(o, data, updatedBy, columns...)
	if err != nil {
		fmt.Println("Error updating tg_user ", err)
	}
	return
}

func updateRowByColumns(o orm.Ormer, data *TgUser, updatedBy string, columns ...string) (err error) {
	_, err = o.Update(data, columns...)
	return
}

func GetUserByVerificationCode(verificationCode string) (v *TgUser, err error) {

	o := orm.NewOrm()
	v = &TgUser{}
	err = o.QueryTable(new(TgUser)).Filter("verificationCode", verificationCode).OrderBy("-id").One(v)
	if err != nil {
		fmt.Println("error GetUserByVerificationCode ", err)
		return
	}

	return
}

func GetUserByEmail(email string) (v *TgUser, err error) {

	o := orm.NewOrm()
	v = &TgUser{}
	err = o.QueryTable(new(TgUser)).Filter("email", email).OrderBy("-id").One(v)
	if err != nil && err == orm.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		fmt.Println("error getting tg_user by email ", err)
		return
	}
	return
}
