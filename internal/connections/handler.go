package connections

import (
	"backuper/internal/db"
	"gorm.io/gorm"
)

func NewConnection(db *gorm.DB, cn db.Connection) db.Connection {
	db.Create(&cn)

	return cn
}
