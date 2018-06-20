package lib

// Config describes the pepster config
type Config struct {
	Token    string `description:"Bot token for authentication"`
	APIKey   string `description:"osu! API key"`
	Database string `description:"Path to sqlite database file"`
}

var DefaultCfg = Config{
	Database: "pepster.db",
}
