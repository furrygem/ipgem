package repository

import (
	"database/sql"

	"github.com/furrygem/ipgem/api/internal/logger"
	"github.com/furrygem/ipgem/api/internal/models"
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

func (sqliterepo *SQLiteRepository) List() (error, *models.RecordList) {
	// FIXME: Explicit field names in select needed
	rows, err := sqliterepo.dbConn.Query("SELECT * FROM records")
	l := logger.GetLogger()
	l.Info(rows)
	if err != nil {
		return err, nil
	}
	var dest models.RecordList = models.RecordList{}
	for rows.Next() {
		record := models.Record{}
		// BUG: Doesn't read the date time fields correctly: "created_at": "0001-01-01T00:00:00Z",
		// STYLE: improve the style here, line is too long
		rows.Scan(&record.RecordID, &record.DomainName, &record.RecordType, &record.Value, &record.TTL, &record.CreatedAt, &record.UpdatedAt)
		dest = append(dest, record)
	}
	if err != nil {
		return err, nil
	}
	return nil, &dest
}

func (sqliterepo *SQLiteRepository) Update() error {
	return nil
}

func (sqliterepo *SQLiteRepository) Retrieve() error {
	return nil
}

func (sqliterepo *SQLiteRepository) Delete() error {
	return nil
}
