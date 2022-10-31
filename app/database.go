package app

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"idcra-telegram-scheduler/helper"
	"log"
	"time"
)

func NewDB() *gorm.DB {
	dbUname := helper.Getenv("DATABASE_USERNAME", "")
	dbPass := helper.Getenv("DATABASE_PASSWORD", "")

	dbUrl := helper.Getenv("DATABASE_HOST", "")
	dbPort := helper.Getenv("DATABASE_PORT", "")
	dbName := helper.Getenv("DATABASE_NAME", "")

	log.Printf("Connection to Database....")
	state := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", dbUname, dbPass, dbUrl, dbPort, dbName)
	db, err := gorm.Open(mysql.Open(state), &gorm.Config{})
	helper.PanicIfError(err)

	sqlDB, err := db.DB()
	helper.PanicIfError(err)

	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetConnMaxLifetime(60 * time.Minute)

	log.Printf("Connected...")
	return db
}
