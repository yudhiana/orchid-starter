package config

type LocalConfig struct {
	DatabaseDebug      bool   `env:"DATABASE_DEBUG" envDefault:"false"`
	ElasticsearchDebug bool   `env:"ES_DEBUG" envDefault:"false"`
	RedisDebug         bool   `env:"REDIS_DEBUG" envDefault:"false"`
	SentryDsn          string `env:"SENTRY_DSN"`
	LogLevel           string `env:"LOG_LEVEL" envDefault:"INFO"`

	AppName    string `env:"APP_NAME" envDefault:"orchid-starter"`
	AppPort    string `env:"APP_PORT" envDefault:"8080"`
	AppHost    string `env:"APP_HOST" envDefault:"0.0.0.0"`
	AppVersion string `env:"APP_VERSION" envDefault:"1.0.0"`
	AppEnv     string `env:"APP_ENV"`

	// mysql config
	MySQLConfig MySQLConfig

	// elasticsearch config
	EsConfig EsConfig

	// redis config
	RedisConfig RedisConfig

	// logger config
	LoggerConfig LoggerConfig
}

type EsConfig struct {
	ESAddresses           string `env:"ES_ADDRESSES,required"`
	ESIdleTimeOut         int    `env:"ES_IDLE_TIMEOUT" envDefault:"60"`
	ESMaxIdleConns        int    `env:"ES_MAX_IDLE_CONNS" envDefault:"100"`
	ESMaxIdleConnsPerHost int    `env:"ES_MAX_IDLE_CONN_PER_HOST" envDefault:"10"`
	ESMaxConnsPerHost     int    `env:"ES_MAX_CONNS_PER_HOST" envDefault:"100"`
}

type MySQLConfig struct {
	MySQLHost         string `env:"DATABASE_HOST,required"`
	MySQLPort         string `env:"DATABASE_PORT,required"`
	MySQLDatabaseName string `env:"DATABASE_NAME,required"`
	MySQLUsername     string `env:"USERNAME_DB,required"`
	MySQLPassword     string `env:"PASSWORD_DB,required"`

	MySQLSetMaxIdleConns      int `env:"MAX_IDLE_CONNS_MYSQL" envDefault:"5"`
	MySQLSetMaxOpenConns      int `env:"MAX_OPEN_CONNS_MYSQL" envDefault:"10"`
	MySQLSetMaxConnLifetime   int `env:"CONN_MAX_LIFETIME_MYSQL" envDefault:"60"`
	MySQLSetMaxIdleConnection int `env:"MAX_IDLE_CONNECTION_MYSQL" envDefault:"5"`
}

type RedisConfig struct {
	RedisHost     string `env:"REDIS_HOST" envDefault:"localhost"`
	RedisPort     string `env:"REDIS_PORT" envDefault:"6379"`
	RedisUsername string `env:"REDIS_USERNAME"`
	RedisPassword string `env:"REDIS_PASSWORD"`
	RedisDB       int    `env:"REDIS_DB"`

	RedisPoolSize        int `env:"REDIS_POOL_SIZE" envDefault:"10"`
	RedisMinIdleConn     int `env:"REDIS_MIN_IDLE_CONN" envDefault:"2"`
	RedisConnMaxIdleTime int `env:"REDIS_CONN_MAX_IDLE_TIME" envDefault:"600"` // seconds (10 min)
}

type LoggerConfig struct {
	LoggerFileLocation    string `env:"LOGGER_FILE_LOCATION"`
	LoggerFileTdrLocation string `env:"LOGGER_TDR_FILE_LOCATION"`
	LoggerFileMaxAge      int    `env:"LOGGER_FILE_MAX_AGE"`
	LoggerStdout          bool   `env:"LOGGER_STDOUT"`
}
