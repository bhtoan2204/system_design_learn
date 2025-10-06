package settings

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

// type SecurityConfig struct {
// 	JWTAccessSecret      string `mapstructure:"jwt_access_secret"`
// 	JWTRefreshSecret     string `mapstructure:"jwt_refresh_secret"`
// 	JWTAccessExpiration  string `mapstructure:"jwt_access_expiration"`
// 	JWTRefreshExpiration string `mapstructure:"jwt_refresh_expiration"`
// 	HMACSecret           string `mapstructure:"hmac_secret"`
// }

type LogConfig struct {
	LogLevel   string `mapstructure:"log_level"`
	FilePath   string `mapstructure:"file_path"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}

type RedisConfig struct {
	ConnectionURL       string   `mapstructure:"connection_url"`
	Password            string   `mapstructure:"password"`
	DB                  int      `mapstructure:"db"`
	UseSentinel         bool     `mapstructure:"use_sentinel"`
	SentinelMasterName  string   `mapstructure:"sentinel_master_name"`
	SentinelServers     []string `mapstructure:"sentinel_servers"`
	PoolSize            int      `mapstructure:"pool_size"`
	DialTimeoutSeconds  int      `mapstructure:"dial_timeout_seconds"`
	ReadTimeoutSeconds  int      `mapstructure:"read_timeout_seconds"`
	WriteTimeoutSeconds int      `mapstructure:"write_timeout_seconds"`
	IdleTimeoutSeconds  int      `mapstructure:"idle_timeout_seconds"`
	MaxIdleConn         int      `mapstructure:"max_idle_conn_number"`
	MaxActiveConn       int      `mapstructure:"max_active_conn_number"`
	SentinelPassword    string   `mapstructure:"sentinel_password"`
}

type ServiceConfig struct {
	PaymentServiceUrl string `mapstructure:"payment_service_url"`
}

type Config struct {
	Server    ServerConfig `mapstructure:"server"`
	LogConfig LogConfig    `mapstructure:"log"`
	// SecurityConfig SecurityConfig `mapstructure:"security"`
	RedisConfig RedisConfig   `mapstructure:"redis"`
	Service     ServiceConfig `mapstructure:"service"`
}
