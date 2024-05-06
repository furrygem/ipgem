package service

import (
	"github.com/furrygem/ipgem/api/internal/logger"
	"github.com/furrygem/ipgem/api/internal/models"
	"github.com/furrygem/ipgem/api/internal/repository"
	"github.com/google/uuid"
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

func (dnscrud *DNSCrud) AddRecord(recordDTO *models.RecordDTO) (models.Record, error) {
	record := &models.Record{}
	record.DomainName = recordDTO.DomainName
	record.RecordType = recordDTO.RecordType
	record.TTL = recordDTO.TTL
	record.Value = recordDTO.Value
	record.RecordID = uuid.New()
	err, newRecord := dnscrud.Repository.Insert(record)
	if err != nil {
		return newRecord, err
	}
	return newRecord, nil
}

func (dnscrud *DNSCrud) UpdateRecord(id string, recordDTO *models.RecordDTO) (models.Record, error) {
	record := models.Record{}
	record.DomainName = recordDTO.DomainName
	record.RecordType = recordDTO.RecordType
	record.TTL = recordDTO.TTL
	record.Value = recordDTO.Value
	err, record := dnscrud.Repository.Update(id, &record)
	if err != nil {
		return record, err
	}
	return record, nil
}

func (dnscrud *DNSCrud) DeleteRecord(id string) error {
	err := dnscrud.Repository.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
