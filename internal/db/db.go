package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewSqlite(path string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// Migrate the schema
	err = db.AutoMigrate(&Connection{})
	if err != nil {
		panic(err)
	}

	return db
}
