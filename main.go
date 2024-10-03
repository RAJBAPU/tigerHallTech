package main

import (
	"fmt"
	"log"
	"simpl_pr/handler"
	"simpl_pr/middleware"
	"simpl_pr/service"

	"simpl_pr/persistence"

	"github.com/astaxie/beego/orm"
	_ "github.com/astaxie/beego/orm" // Import Beego ORM package
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql" // Import MySQL driver package
)

func main() {
	fmt.Println("Starting")
	err := registerDB()
	if err != nil {
		log.Fatalf("<MysqlDB> Error in kbyp  db registration: %v", err)
	}
	go service.ProcessMessages()
	router := gin.Default()

	dbRunMode := "mysql"
	tgConfig := persistence.NewTgConfigPersistence(dbRunMode)
	tgUser := persistence.NewTgUsergPersistence(dbRunMode)
	tgTigerDetails := persistence.NewTgTigerDetailsPersistence(dbRunMode)
	tgTigerImages := persistence.NewTgTigerImagesPersistence(dbRunMode)
	tgDependency := &service.Tiger{
		Configs:      tgConfig,
		User:         tgUser,
		TigerDetails: tgTigerDetails,
		TigerImages:  tgTigerImages,
	}

	tgUserDependency := &service.User{
		Configs: tgConfig,
		User:    tgUser,
	}

	customer := &middleware.CustomerSvc{
		User:    tgUser,
		Configs: tgConfig,
	}

	router.POST("/signup", handler.SignupPage(tgUserDependency))
	router.GET("/verifyemail/:verificationCode", handler.VerifyEmail(tgUserDependency))
	router.POST("/login", handler.SignInUser(tgUserDependency))
	router.GET("/me", middleware.DeserializeUser(customer), handler.GetUserDetails(tgUserDependency))

	router.POST("/tigers", middleware.DeserializeUser(customer), handler.PostTigerDetails(tgDependency))
	router.GET("/tigers", middleware.DeserializeUser(customer), handler.GetAllTigers(tgDependency))
	router.POST("/tigers/sightings", middleware.DeserializeUser(customer), handler.PostSightingDetails(tgDependency))
	router.GET("/tigers/sightings", middleware.DeserializeUser(customer), handler.GetSightingDetails(tgDependency))

	router.Run("localhost:8080")
}

func registerDB() error {
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

/// move many var to yp_config
