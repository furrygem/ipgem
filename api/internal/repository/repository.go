package repository

import "github.com/furrygem/ipgem/api/internal/models"

// type Filter interface{}
// type Identifier interface{}
// type Record interface{}

type Repository interface {
	Open() error
	Close() error
	List() (error, *models.RecordList)
	Update(string, *models.Record) (error, models.Record)
	Retrieve(string) (error, models.Record)
	Insert(*models.Record) (error, models.Record)
	Delete() error
}
