package mysql

import (
	"log"
	"os"
	"sync"
	"time"

	"orchid-starter/config"
	"orchid-starter/internal/common"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/yudhiana/logos"
)

var mysqlInstance *gorm.DB
var mysqlOnce sync.Once

func getLogger(debug bool) logger.Interface {
	if !debug {
		return logger.Default
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: false,       // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)

	return newLogger
}

// GetDBConnection gets DB connection
func GetMySQLConnection(cfg *config.LocalConfig) *gorm.DB {
	mysqlOnce.Do(func() {

		tmpl := "{{username}}:{{password}}@tcp({{host}}:{{port}})/{{db_name}}?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=True&loc=Local"

		//setup connection to database
		URI, err := common.Render(tmpl, map[string]any{
			"username": cfg.MySQLConfig.MySQLUsername,
			"password": cfg.MySQLConfig.MySQLPassword,
			"host":     cfg.MySQLConfig.MySQLHost,
			"port":     cfg.MySQLConfig.MySQLPort,
			"db_name":  cfg.MySQLConfig.MySQLDatabaseName,
		})

		if err != nil {
			panic(err)
		}

		logos.NewLogger().Info("Initialize MySQL connection", "URI", URI)
		db, err := gorm.Open(mysql.New(mysql.Config{
			DSN: URI,
		}), gormConfig(cfg.DatabaseDebug))
		if err != nil {
			panic(err)
		}

		sqlDB, err := db.DB()
		if err != nil {
			panic(err)
		}

		sqlDB.SetConnMaxIdleTime(time.Second * time.Duration(cfg.MySQLConfig.MySQLMaxIdleConnection))
		sqlDB.SetConnMaxLifetime(time.Second * time.Duration(cfg.MySQLConfig.MySQLMaxConnLifetime))
		sqlDB.SetMaxIdleConns(cfg.MySQLConfig.MySQLMaxIdleConns)
		sqlDB.SetMaxOpenConns(cfg.MySQLConfig.MySQLMaxOpenConns)

		mysqlInstance = db

	})

	return mysqlInstance
}

func gormConfig(debug bool) *gorm.Config {
	return &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: getLogger(debug),
	}
}

func GetMockSQLConnection() (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), gormConfig(true))
	if err != nil {
		panic(err)
	}

	return gormDB, mock
}
