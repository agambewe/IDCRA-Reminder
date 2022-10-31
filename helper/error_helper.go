package helper

import (
	"gorm.io/gorm"
	"log"
)

func PanicIfError(err error) {
	if err != nil {
		log.Panic(err.Error())
	}
}

func RecoveryIfPanic(db *gorm.DB) {

	log.Printf("Recovery called...")

	if r := recover(); r != nil {
		log.Printf("Recovery: %s", r)
		db.Rollback()

		sqlDb, _ := db.DB()
		sqlDb.Close()

		db.Commit()
	}

	log.Printf("Panic handled...")
}
