package config

type DBConnection struct {
	Name     string `mapstructure:"name"`
	Type     string `mapstructure:"type"` // "postgres" or "mysql"
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

type Config struct {
	Connections []DBConnection `mapstructure:"connections"`
	Schedule    string         `mapstructure:"schedule"`
	Telegram    struct {
		Token  string `mapstructure:"token"`
		ChatID string `mapstructure:"chat_id"`
	} `mapstructure:"telegram"`
}
