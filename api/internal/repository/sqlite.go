package repository

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// List(Filter) (error, interface{})
// Update(Identifier, Record) (error, interface{})
// Retrieve(Identifier) (error, interface{})
// Delete(Identifier) error

type SQLiteRepository struct {
	dbPath string
	dbConn *sql.DB
}

func NewSQLiteRepository(dbPath string) (*SQLiteRepository, error) {
	return &SQLiteRepository{
		dbPath: dbPath,
	}, nil
}

func (sqliterepo *SQLiteRepository) Open() error {
	conn, err := sql.Open("sqlite3", sqliterepo.dbPath)
	if err != nil {
		return err
	}
	sqliterepo.dbConn = conn
	return nil
}

func (sqliterepo *SQLiteRepository) Close() error {
	return sqliterepo.dbConn.Close()
}
