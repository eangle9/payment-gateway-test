package foundation

import (
	"context"
	"fmt"
	"pg/platform/hlog"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // Use PostgreSQL instead of CockroachDB
	_ "github.com/golang-migrate/migrate/v4/source/file"       // Keep the file source
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func InitDB(url string, log hlog.Logger) *pgxpool.Pool {
	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Fatal(context.Background(), fmt.Sprintf("Failed to connect to database: %v", err))
	}

	idleConnTimeout := viper.GetDuration("database.idle_conn_timeout")
	if idleConnTimeout == 0 {
		idleConnTimeout = 4 * time.Minute
	}

	config.ConnConfig.Logger = log.Named("pgx")
	config.MaxConnIdleTime = idleConnTimeout
	conn, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Fatal(context.Background(), fmt.Sprintf("Failed to connect to database: %v", err))
	}

	// Test the connection (Replace with Ping or proper SQL)
	if _, err := conn.Exec(context.Background(),
		"SELECT table_name FROM information_schema.tables WHERE table_schema = 'public';"); err != nil {
		log.Fatal(context.Background(), fmt.Sprintf("Failed to ping database: %v", err))
	}

	return conn
}

func InitiateMigration(path, conn string, log hlog.Logger) *migrate.Migrate {
	conn = fmt.Sprintf("postgres://%s", strings.Split(conn, "://")[1]) // Change to postgres
	m, err := migrate.New(fmt.Sprintf("file://%s", path), conn)
	if err != nil {
		log.Fatal(context.Background(), "could not create migrator", zap.Error(err))
	}
	return m
}

func UpMigration(m *migrate.Migrate, log hlog.Logger) {
	err := m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal(context.Background(), "could not migrate", zap.Error(err))
	}
}
