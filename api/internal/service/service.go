package service

import "github.com/furrygem/ipgem/api/internal/repository"

type DNSCrud struct {
	repository repository.Repository
}

func NewService(repo repository.Repository) *DNSCrud {
	return &DNSCrud{
		repository: repo,
	}
}

func (dnscrud *DNSCrud) ListRecords() {}

func (dnscrud *DNSCrud) AddRecord() {}

func (dnscrud *DNSCrud) UpdateRecord() {}

func (dnscrud *DNSCrud) DeleteRecord() {}
