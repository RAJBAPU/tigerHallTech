package initializers

import "github.com/astaxie/beego/orm"

func RegisterDB() error {
	dbUser := "root"
	dbPass := "Kreditbee@123"
	dbHost := "127.0.0.1"
	dbPort := "3306"
	dbName := "sys"

	// Register the MySQL driver
	err := orm.RegisterDriver("mysql", orm.DRMySQL)
	if err != nil {
		return err
	}
	// Set the database connection string
	dataSource := dbUser + ":" + dbPass + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8"

	// Register the database connection
	err = orm.RegisterDataBase("default", "mysql", dataSource)
	if err != nil {
		return err
	}
	return nil
}
