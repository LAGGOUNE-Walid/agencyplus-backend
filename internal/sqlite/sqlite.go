package sqlite

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/sqlite3"
	_ "github.com/mattn/go-sqlite3"
)

type SQLite struct {
	db *sql.DB
}

func New(mode string, path string) (*SQLite, error) {
	var dataSource string
	if mode == "memory" {
		dataSource = ":memory:"
	} else {
		dataSource = path
	}
	db, err := sql.Open("sqlite3", dataSource)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s db: %w", mode, err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping %s db: %w", mode, err)
	}
	return &SQLite{db: db}, nil
}

func NewFromFile(path string) (*SQLite, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to sqlite: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping failed: %w", err)
	}
	return &SQLite{db: db}, nil
}

func NewInMemory(t *testing.T) (*SQLite, error) {
	db, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		return nil, fmt.Errorf("failed to open in-memory DB: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping in-memory DB: %w", err)
	}
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		t.Fatal("failed to create database instance:", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"sqlite3", driver)
	if err != nil {
		t.Fatal("failed to initialize migrate:", err)
	}

	m.Down()
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		t.Fatal("failed to apply migrations:", err)
	}
	return &SQLite{db: db}, nil
}

func (s *SQLite) GetDB() *sql.DB {
	return s.db
}

func (s *SQLite) Close() error {
	return s.db.Close()
}
