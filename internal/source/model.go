package source

import "gorm.io/gorm"

type Source struct {
	gorm.Model
	Title    string
	Link     string
	Host     string
	Port     int
	DB       string
	Username string
	Password string
	DBType   string
}
