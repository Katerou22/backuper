package backuper

import (
	"fmt"
)

type DBSource interface {
	GetTitle() string
	GetDBType() string
	GetHost() string
	GetPort() int
	GetUsername() string
	GetPassword() string
	GetDBName() string
}
type Backuper interface {
	Backup() (string, error)
}

func NewBackuper(src DBSource) (Backuper, error) {
	switch src.GetDBType() {
	case "mysql":
		return &MySQLBackuper{src}, nil
	case "postgres":
		return &PostgresBackuper{src}, nil
	default:
		return nil, fmt.Errorf("unsupported db type: %s", src.GetDBType())
	}
}
