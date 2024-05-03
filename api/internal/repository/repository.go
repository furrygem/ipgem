package repository

import "github.com/furrygem/ipgem/api/internal/models"

// type Filter interface{}
// type Identifier interface{}
// type Record interface{}

type Repository interface {
	Open() error
	Close() error
	List() (error, *models.RecordList)
	Update() error
	Retrieve() error
	Delete() error
}
