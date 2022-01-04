package storage

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

var SESSION *gorm.DB

func Load(sessionName string) {
	db, err := gorm.Open(sqlite.Open(sessionName), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Panicln(err)
	}
	SESSION = db
	dB, _ := db.DB()
	dB.SetMaxOpenConns(100)
	// db.DB().SetMaxOpenConns(100)

	// Create tables if they don't exist
	SESSION.AutoMigrate(&Session{}, &Peer{})

}
