package utils

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Up runs migrations
func Up(ctx context.Context, connStr string) error {
	// migrationsDir "file:///migrations"
	tmpDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("%v\n", err)
	}
	migrationsDir := "file://" + tmpDir + "/migrations"

	m, err := migrate.New(migrationsDir, connStr)
	if err != nil {
		log.Panic("cannot create new migration", migrationsDir, err)
	}
	m.Log = &customLog{verbose: true}

	err = m.Up()
	if errors.Is(err, migrate.ErrNoChange) || errors.Is(err, migrate.ErrNilVersion) {
		log.Println("no migration:", err)
		return nil
	}

	if err != nil {
		log.Println("migration error:", err)
	}

	return err
}

// customLog represents the logger
type customLog struct {
	verbose bool
}

// Printf prints out formatted string into a log
func (l *customLog) Printf(format string, v ...interface{}) {
	if l.verbose {
		log.Printf(format, v...)
	} else {
		fmt.Fprintf(os.Stderr, format, v...)
	}
}

// Println prints out args into a log
func (l *customLog) Println(args ...interface{}) {
	if l.verbose {
		log.Println(args...)
	} else {
		fmt.Fprintln(os.Stderr, args...)
	}
}

// Verbose shows if verbose print enabled
func (l *customLog) Verbose() bool {
	return l.verbose
}
