package pepster

// Config describes the pepster config
type Config struct {
	Token      string `description:"Bot token for authentication"`
	APIKey     string `description:"osu! API key"`
	LogChannel string `description:"Channel to log debug information."`

	Database     string `description:"Path to sqlite database file"`
	RedisAddress string `description:"Redis connection address ('host:port')"`
	RedisDB      int    `description:"Redis database to connect to"`
}

var DefaultCfg = Config{
	Database: "pepster.db",
}
