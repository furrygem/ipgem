package service

import (
	"github.com/furrygem/ipgem/api/internal/logger"
	"github.com/furrygem/ipgem/api/internal/models"
	"github.com/furrygem/ipgem/api/internal/repository"
)

type DNSCrud struct {
	Repository repository.Repository
}

func NewService(repo repository.Repository) *DNSCrud {
	err := repo.Open()
	l := logger.GetLogger()
	if err != nil {
		l.Error(err)
	}
	return &DNSCrud{
		Repository: repo,
	}
}

func (dnscrud *DNSCrud) CloseConn() {
	dnscrud.Repository.Close()
}

func (dnscrud *DNSCrud) ListRecords() (*models.RecordList, error) {
	err, records := dnscrud.Repository.List()
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (dnscrud *DNSCrud) RetrieveRecord(id string) (models.Record, error) {
	err, record := dnscrud.Repository.Retrieve(id)
	if err != nil {
		return record, err
	}
	return record, nil
}

func (dnscrud *DNSCrud) AddRecord() {}

func (dnscrud *DNSCrud) UpdateRecord() {}

func (dnscrud *DNSCrud) DeleteRecord() {}
