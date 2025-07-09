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

// Implementing DBSource interface

func (s *Source) GetTitle() string    { return s.Title }
func (s *Source) GetDBType() string   { return s.DBType }
func (s *Source) GetHost() string     { return s.Host }
func (s *Source) GetPort() int        { return s.Port }
func (s *Source) GetUsername() string { return s.Username }
func (s *Source) GetPassword() string { return s.Password }
func (s *Source) GetDBName() string   { return s.DB }
