package repository

import (
	"database/sql"
	"time"

	"github.com/furrygem/ipgem/api/internal/logger"
	"github.com/furrygem/ipgem/api/internal/models"
	_ "github.com/mattn/go-sqlite3"
)

// List(Filter) (error, interface{})
// Update(Identifier, Record) (error, interface{})
// Retrieve(Identifier) (error, interface{})
// Delete(Identifier) error

type SQLiteRepository struct {
	dbPath                  string
	dbConn                  *sql.DB
	listRecordsStatement    *sql.Stmt
	retrieveRecordStatement *sql.Stmt
	updateRecordStatement   *sql.Stmt
}

const listRecordsQuery = `SELECT
	record_id,
	domain_name,
	record_type,
	value,
	ttl,
	CAST(created_at AS INTEGER),
	CAST(updated_at AS INTEGER)
FROM records`

const retrieveRecordQuery = `SELECT
	record_id,
	domain_name,
	record_type,
	value,
	ttl,
	CAST(created_at AS INTEGER),
	CAST(updated_at AS INTEGER)
FROM records
WHERE record_id=$1`

const updateRecordQuery = `UPDATE records
SET
	domain_name = $1,
	record_type = $2,
	value = $3,
	ttl = $4,
	updated_at = unixepoch('now')
WHERE record_id = $5
RETURNING
	record_id,
	domain_name,
	record_type,
	value,
	ttl,
	CAST(created_at AS INTEGER),
	CAST(updated_at AS INTEGER)`

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
	sqliterepo.listRecordsStatement, err = conn.Prepare(listRecordsQuery)
	if err != nil {
		return err
	}
	sqliterepo.retrieveRecordStatement, err = conn.Prepare(retrieveRecordQuery)
	if err != nil {
		return err
	}
	sqliterepo.updateRecordStatement, err = conn.Prepare(updateRecordQuery)
	if err != nil {
		return err
	}
	return nil
}

func (sqliterepo *SQLiteRepository) Close() error {
	return sqliterepo.dbConn.Close()
}

func (sqliterepo *SQLiteRepository) List() (error, *models.RecordList) {
	rows, err := sqliterepo.listRecordsStatement.Query()
	l := logger.GetLogger()
	if err != nil {
		return err, nil
	}
	var dest models.RecordList = models.RecordList{}
	for rows.Next() {
		record := models.Record{}
		// BUG: Doesn't read the date time fields correctly: "created_at": "0001-01-01T00:00:00Z",
		// STYLE: improve the style here, line is too long
		var createdAtTs int64
		var updatedAtTs int64
		err := rows.Scan(&record.RecordID,
			&record.DomainName,
			&record.RecordType,
			&record.Value,
			&record.TTL,
			&createdAtTs,
			&updatedAtTs)
		if err != nil {
			return err, nil
		}
		record.CreatedAt = time.Unix(createdAtTs, 0)
		record.UpdatedAt = time.Unix(updatedAtTs, 0)
		dest = append(dest, record)
	}
	if err != nil {
		return err, nil
	}
	return nil, &dest
}

func (sqliterepo *SQLiteRepository) Update(id string, new *models.Record) (error, models.Record) {
	row := sqliterepo.updateRecordStatement.QueryRow(new.DomainName, new.RecordType, new.Value, new.TTL, id)
	record := models.Record{}
	var createdAtTs int64
	var updatedAtTs int64
	err := row.Scan(&record.RecordID,
		&record.DomainName,
		&record.RecordType,
		&record.Value,
		&record.TTL,
		&createdAtTs,
		&updatedAtTs)
	if err != nil {
		return err, record
	}
	record.CreatedAt = time.Unix(createdAtTs, 0)
	record.UpdatedAt = time.Unix(updatedAtTs, 0)
	return nil, record
}

func (sqliterepo *SQLiteRepository) Retrieve(id string) (error, models.Record) {
	row := sqliterepo.retrieveRecordStatement.QueryRow(id)
	record := models.Record{}
	var createdAtTs int64
	var updatedAtTs int64
	err := row.Scan(&record.RecordID,
		&record.DomainName,
		&record.RecordType,
		&record.Value,
		&record.TTL,
		&createdAtTs,
		&updatedAtTs)
	if err != nil {
		return err, record
	}
	record.CreatedAt = time.Unix(createdAtTs, 0)
	record.UpdatedAt = time.Unix(updatedAtTs, 0)
	return nil, record
}

func (sqliterepo *SQLiteRepository) Delete() error {
	return nil
}
