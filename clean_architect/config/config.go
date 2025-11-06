package config

type Config struct {
	Server   Server
	Database Database
}

type Server struct {
	Port int `env:"PORT" default:"8080"`
}

type Database struct {
	DriverName             string `env:"DB_DRIVER_NAME"`
	ConnectionURL          string `env:"DB_CONNECTION_URL" json:"-"` //zap ignore
	MaxOpenConnNumber      int    `env:"DB_MAX_OPEN_CONN_NUMBER"`
	MaxIdleConnNumber      int    `env:"DB_MAX_IDLE_CONN_NUMBER"`
	ConnMaxLifeTimeSeconds int64  `env:"DB_CONN_MAX_LIFE_TIME_SECONDS"`
	ConnMaxIdleTimeSeconds int64  `env:"DB_CONN_MAX_IDLE_TIME_SECONDS"`
}
