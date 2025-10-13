package migrator

import (
	"database/sql"
	"fmt"

	"transaction/pkg/config"

	migrate "github.com/rubenv/sql-migrate"
)

type Migrator struct {
	dialect       string
	dbConfig      config.DatabaseConfig
	migrationsDir *migrate.FileMigrationSource
}

func New(cfg config.DatabaseConfig) *Migrator {
	return &Migrator{
		dialect: "postgres",
		migrationsDir: &migrate.FileMigrationSource{
			Dir: "migrations",
		},
		dbConfig: cfg,
	}
}

func (m *Migrator) Up() error {
	conn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		m.dbConfig.User,
		m.dbConfig.Password,
		m.dbConfig.Host,
		m.dbConfig.Port,
		m.dbConfig.DBName,
		m.dbConfig.SSLMode,
	)

	db, err := sql.Open(m.dialect, conn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	migrate.SetTable("migrations")

	n, err := migrate.Exec(db, m.dialect, m.migrationsDir, migrate.Up)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	fmt.Printf("Applied %d migrations\n", n)
	return nil
}

func (m *Migrator) Down() error {
	conn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		m.dbConfig.User,
		m.dbConfig.Password,
		m.dbConfig.Host,
		m.dbConfig.Port,
		m.dbConfig.DBName,
		m.dbConfig.SSLMode,
	)

	db, err := sql.Open(m.dialect, conn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	migrate.SetTable("migrations")

	n, err := migrate.Exec(db, m.dialect, m.migrationsDir, migrate.Down)
	if err != nil {
		return fmt.Errorf("failed to rollback migrations: %w", err)
	}

	fmt.Printf("Rolled back %d migrations\n", n)
	return nil
}
