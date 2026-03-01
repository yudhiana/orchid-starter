package config

import "fmt"

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
	MySQLHost         string `env:"MYSQL_DATABASE_HOST,required"`
	MySQLPort         string `env:"MYSQL_DATABASE_PORT,required"`
	MySQLDatabaseName string `env:"MYSQL_DATABASE_NAME,required"`
	MySQLUsername     string `env:"MYSQL_USERNAME,required"`
	MySQLPassword     string `env:"MYSQL_PASSWORD,required"`

	MySQLMaxIdleConns      int `env:"MYSQL_MAX_IDLE_CONNS" envDefault:"5"`
	MySQLMaxOpenConns      int `env:"MYSQL_MAX_OPEN_CONNS" envDefault:"10"`
	MySQLMaxConnLifetime   int `env:"MYSQL_CONN_MAX_LIFETIME" envDefault:"60"`
	MySQLMaxIdleConnection int `env:"MYSQL_MAX_IDLE_CONNECTION" envDefault:"5"`
}

func (m MySQLConfig) DSN() string {
	tmpl := "%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=True&loc=Local"
	return fmt.Sprintf(
		tmpl,
		m.MySQLUsername,
		m.MySQLPassword,
		m.MySQLHost,
		m.MySQLPort,
		m.MySQLDatabaseName,
	)
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

func (rds RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%s", rds.RedisHost, rds.RedisPort)
}

type RabbitMQConfig struct {
	RabbitMQHost     string `env:"RABBITMQ_HOST" envDefault:"rabbitmq"`
	RabbitMQPort     string `env:"RABBITMQ_PORT" envDefault:"5672"`
	RabbitMQUser     string `env:"RABBITMQ_USER" envDefault:"guest"`
	RabbitMQPassword string `env:"RABBITMQ_PASSWORD" envDefault:"guest"`
}

func (r RabbitMQConfig) AmqpURI() string {
	return fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		r.RabbitMQUser,
		r.RabbitMQPassword,
		r.RabbitMQHost,
		r.RabbitMQPort,
	)
}

type LoggerConfig struct {
	LoggerFileLocation    string `env:"LOGGER_FILE_LOCATION"`
	LoggerFileTdrLocation string `env:"LOGGER_TDR_FILE_LOCATION"`
	LoggerFileMaxAge      int    `env:"LOGGER_FILE_MAX_AGE"`
	LoggerStdout          bool   `env:"LOGGER_STDOUT"`
}
