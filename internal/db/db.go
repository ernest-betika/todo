package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var db DB

type SQLOperations interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type DB interface {
	SQLOperations
	Begin() (*sql.Tx, error)
	Close() error
	InTransaction(ctx context.Context, operations func(context.Context, SQLOperations) error) error
	Ping() error
	Valid() bool
}

type RowScanner interface {
	Scan(dest ...any) error
}

type AppDB struct {
	*sql.DB
	valid bool
}

type pgSQLOperations struct {
	*sql.Tx
}

func (h *AppDB) InTx(ctx context.Context, f func(tx *sql.Tx) error) error {

	tx, err := h.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	if err := f(tx); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	if err = tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		// Optionally return an error here to signal the commit failure
		// err = fmt.Errorf("failed to commit transaction: %w", err)
		return err
	}

	return nil
}

func InitDB() DB {
	return initDBWithURL(os.Getenv("DATABASE_URL"))
}

func initDBWithURL(databaseURL string) DB {

	if databaseURL == "" {
		log.Fatal("database url is empty")
	}

	appDB, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalf("sql open error %v", err)
	}

	db = &AppDB{
		DB:    appDB,
		valid: true,
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("database ping error %v", err)
	}

	return db
}

func (db *AppDB) Valid() bool {
	return db.valid
}

//option 1 
//WithTransaction runs the given function f within a transaction (ideally as a utility helper)
//can be called from the business logic assuming each and every table has it's own repository
func WithTransaction(db DB, f func(SQLOperations) error) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	pgSQLOperations := &pgSQLOperations{
		Tx: tx,
	}

	if err := f(pgSQLOperations); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	return tx.Commit()
}

//option 2
//InTransaction runs the given function f within a transaction (ideally if you are passing db connection as a parameter)
func (db *AppDB) InTransaction(ctx context.Context, f func(context.Context, SQLOperations) error) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	pgSQLOperations := &pgSQLOperations{
		Tx: tx,
	}

	if err := f(ctx, pgSQLOperations); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	return tx.Commit()
}

func GetDB() DB {
	return db
}
