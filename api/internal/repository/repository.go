package repository

import "github.com/furrygem/ipgem/api/internal/models"

// type Filter interface{}
// type Identifier interface{}
// type Record interface{}

type Repository interface {
	Open() error
	Close() error
	List() (error, *models.RecordList)
	Update(id string, updatedRecord *models.Record) (error, models.Record)
	Retrieve(id string) (error, models.Record)
	Insert(*models.Record) (error, models.Record)
	Delete(id string) error
}
