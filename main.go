package main

import (
	"log"
	"simpl_pr/handler"
	"simpl_pr/middleware"
	"simpl_pr/service"

	"github.com/astaxie/beego/orm"
	_ "github.com/astaxie/beego/orm" // Import Beego ORM package
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql" // Import MySQL driver package
)

func main() {
	err := registerDB()
	if err != nil {
		log.Fatalf("<MysqlDB> Error in kbyp  db registration: %v", err)
	}
	go service.ProcessMessages()
	router := gin.Default()

	router.POST("/signup", handler.SignupPage)
	router.GET("/verifyemail/:verificationCode", handler.VerifyEmail)
	router.POST("/login", handler.SignInUser)
	
	router.GET("/me", middleware.DeserializeUser(), handler.GetUserDetails)

	router.POST("/tigers", middleware.DeserializeUser(), handler.PostTigerDetails)
	router.GET("/tigers", middleware.DeserializeUser(), handler.GetAllTigers)
	router.POST("/tigers/sightings", middleware.DeserializeUser(), handler.PostSightingDetails)
	router.GET("/tigers/sightings", middleware.DeserializeUser(), handler.GetSightingDetails)

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
