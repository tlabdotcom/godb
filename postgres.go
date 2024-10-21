package godb

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

var (
	postgresDB   *bun.DB
	oncePostgres sync.Once
)

// Initialize the PostgreSQL database connection with configurable options.
func initialPostgresSqlDB() error {
	dsn := viper.GetString("DB_POSTGRESQL_DSN")
	if dsn == "" {
		return fmt.Errorf("database DSN is empty")
	}
	log.Infof("Connecting to database with DSN: %v", dsn)

	// Create a new SQL connection using the provided DSN.
	hsqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	if hsqldb == nil {
		return fmt.Errorf("failed to create database connection")
	}

	// Set connection pool settings from configuration, with sensible defaults.
	setMaxOpenConns := viper.GetInt("MAX_OPEN_CONNS")
	if setMaxOpenConns <= 0 {
		setMaxOpenConns = 20 // Default max open connections
	}
	setMaxIdleConns := viper.GetInt("MAX_IDLE_CONNS")
	if setMaxIdleConns <= 0 {
		setMaxIdleConns = 10 // Default max idle connections
	}
	connMaxLifetime := viper.GetDuration("CONN_MAX_LIFETIME")
	if connMaxLifetime == 0 {
		connMaxLifetime = time.Hour // Default connection lifetime
	}

	hsqldb.SetMaxOpenConns(setMaxOpenConns)
	hsqldb.SetMaxIdleConns(setMaxIdleConns)
	hsqldb.SetConnMaxLifetime(connMaxLifetime)

	// Test the database connection with a timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := hsqldb.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Initialize Bun with the SQL connection and PostgreSQL dialect.
	postgresDB = bun.NewDB(hsqldb, pgdialect.New())

	// Add query debugging hooks based on environment variable or config.
	if viper.GetBool("ENABLE_QUERY_DEBUG") {
		postgresDB.AddQueryHook(bundebug.NewQueryHook(
			bundebug.WithVerbose(true),
			bundebug.FromEnv("BUNDEBUG"),
		))
	}
	log.Info("Successfully connected to the database")
	return nil
}

// Get the PostgreSQL database connection, initializing it if necessary.
func GetPostgresDB() *bun.DB {
	oncePostgres.Do(func() {
		err := initialPostgresSqlDB()
		if err != nil {
			log.Errorf("Failed to initialize database: %v", err)
		}
	})
	if postgresDB == nil {
		log.Fatal("Database is not initialized")
	}
	return postgresDB
}

// CloseDatabase gracefully closes the database connection on application shutdown.
func CloseDatabase() error {
	if postgresDB != nil {
		log.Info("Closing database connection")
		return postgresDB.Close()
	}
	return nil
}
