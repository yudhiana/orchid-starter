package bootstrap

import (
	"context"
	"fmt"
	"time"

	"orchid-starter/config"
	"orchid-starter/infrastructure/elastic"
	"orchid-starter/infrastructure/mysql"
	"orchid-starter/infrastructure/redis"
	"orchid-starter/internal/clients"
	"orchid-starter/internal/common"

	"github.com/elastic/go-elasticsearch/v9"
	redisV9 "github.com/redis/go-redis/v9"
	"github.com/yudhiana/logos"
	"gorm.io/gorm"
)

type DirectInjection struct {
	MySQL  *gorm.DB
	ES     *elasticsearch.Client
	Redis  *redisV9.Client
	Client *clients.Client
	Log    *logos.LogEntry
}

// NewDirectInjection creates a new dependency injection container with proper error handling
func NewDirectInjection(cfg *config.LocalConfig) (*DirectInjection, error) {
	logger := logos.NewLogger()
	logger.Info("Initializing dependency injection container...")

	var (
		mysqlDB     *gorm.DB
		esClient    *elasticsearch.Client
		redisClient *redisV9.Client
	)

	// Initialize MySQL connection
	if common.GetBoolEnv("USE_MOCK_CONNECTION", true) {
		db, _ := mysql.GetMockSQLConnection()
		mysqlDB = db
	} else {
		mysqlDB = mysql.GetMySQLConnection(cfg)
		if mysqlDB == nil {
			return nil, fmt.Errorf("failed to initialize MySQL connection")
		}
	}
	logger.Info("MySQL connection established")

	// Initialize Elasticsearch connection
	if common.GetBoolEnv("USE_MOCK_CONNECTION", true) {
		esClient = elastic.GetESMockConnection()
	} else {
		esClient = elastic.NewESConnection(cfg)
		if esClient == nil {
			return nil, fmt.Errorf("failed to initialize Elasticsearch connection")
		}
	}

	if common.GetBoolEnv("USE_MOCK_CONNECTION", true) {
		redisClient, _ = redis.GetMockRedisConnection()
	} else {
		redisClient = redis.GetRedisClient(cfg)
		if redisClient == nil {
			return nil, fmt.Errorf("failed to initialize Redis connection")
		}
	}

	logger.Info("Elasticsearch connection established")

	// Test connections
	if err := testConnections(mysqlDB, redisClient, esClient); err != nil {
		return nil, fmt.Errorf("connection test failed: %w", err)
	}

	di := &DirectInjection{
		MySQL:  mysqlDB,
		ES:     esClient,
		Redis:  redisClient,
		Client: clients.NewClient(),
		Log:    logos.NewLogger(),
	}

	logger.Info("Dependency injection container initialized successfully")
	return di, nil
}

// testConnections verifies that all connections are working
func testConnections(mysqlDB *gorm.DB, redisClient *redisV9.Client, esClient *elasticsearch.Client) error {
	logger := logos.NewLogger()

	// Test MySQL connection
	if sqlDB, err := mysqlDB.DB(); err != nil {
		return fmt.Errorf("failed to get MySQL database instance: %w", err)
	} else if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("MySQL ping failed: %w", err)
	}

	// Test Redis connection
	if redisClient == nil {
		return fmt.Errorf("redis client is nil")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis ping failed: %w", err)
	}

	// Test Elasticsearch connection
	if _, err := esClient.Info(); err != nil {
		return fmt.Errorf("elasticsearch info request failed: %w", err)
	}

	logger.Info("All database connections tested successfully")
	return nil
}

// Close gracefully closes all connections
func (di *DirectInjection) Close() error {
	logger := logos.NewLogger()
	logger.Info("🔴Closing dependency injection resources...")

	var errors []error

	// Close MySQL connection
	if di.MySQL != nil {
		if sqlDB, err := di.MySQL.DB(); err != nil {
			errors = append(errors, fmt.Errorf("failed to get MySQL DB for closing: %w", err))
		} else if err := sqlDB.Close(); err != nil {
			errors = append(errors, fmt.Errorf("failed to close MySQL connection: %w", err))
		} else {
			logger.Info("🔴MySQL connection closed")
		}
	}

	// Close Redis connection

	if di.Redis != nil {
		if err := di.Redis.Close(); err != nil {
			errors = append(errors, fmt.Errorf("failed to close redis connection: %w", err))
		} else {
			logger.Info("🔴Redis connection closed")
		}
	}

	// Elasticsearch client doesn't need explicit closing in most cases
	// but we can add custom cleanup logic here if needed
	logger.Info("🔴Elasticsearch client cleanup completed")

	if len(errors) > 0 {
		return fmt.Errorf("DI cleanup completed with errors: %v", errors)
	}

	logger.Info("🔴Dependency injection resources closed successfully")
	return nil
}

// GetMySQL returns the MySQL database connection
func (di *DirectInjection) GetMySQL() *gorm.DB {
	return di.MySQL
}

// GetElasticsearch returns the Elasticsearch client
func (di *DirectInjection) GetElasticsearch() *elasticsearch.Client {
	return di.ES
}

// GetClient returns the comprehensive client for all domains
func (di *DirectInjection) GetClient() *clients.Client {
	return di.Client
}
