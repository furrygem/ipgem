package repository

type Filter interface{}
type Identifier interface{}
type Record interface{}

type Repository interface {
	Open() error
	Close() error
	List(Filter) (error, interface{})
	Update(Identifier, Record) (error, interface{})
	Retrieve(Identifier) (error, interface{})
	Delete(Identifier) error
}
