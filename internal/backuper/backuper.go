package backuper

import (
	"backuper/internal/source"
	"fmt"
)

type Backuper interface {
	Backup() error
}

func NewBackuper(src *source.Source) (Backuper, error) {
	switch src.DBType {
	case "mysql":
		return &MySQLBackuper{src}, nil
	case "postgres":
		return &PostgresBackuper{src}, nil
	default:
		return nil, fmt.Errorf("unsupported db type: %s", src.DBType)
	}
}
