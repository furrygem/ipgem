package models

import (
	"time"

	"github.com/google/uuid"
)

type Record struct {
	RecordID   uuid.UUID     `json:"record_id"`
	DomainName string        `json:"domain_name"`
	RecordType string        `json:"record_type"`
	Value      string        `json:"value"`
	TTL        time.Duration `json:"ttl"`
	CreatedAt  time.Time     `json:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at"`
}

type RecordList []Record
