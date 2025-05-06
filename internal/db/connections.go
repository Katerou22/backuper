package db

import "gorm.io/gorm"

type Connection struct {
	gorm.Model
	host     string
	port     int
	db       string
	username string
	password string
	dbType   string
}
