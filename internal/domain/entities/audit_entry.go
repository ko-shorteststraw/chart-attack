package entities

import (
	"time"

	"github.com/google/uuid"
)

type AuditEntry struct {
	Id            uuid.UUID
	CreatedAt     time.Time
	UserId        uuid.UUID
	PatientId     *uuid.UUID
	Action        string
	EntityType    string
	EntityId      uuid.UUID
	FieldsChanged string
	IPAddress     string
	UserAgent     string
}

func NewAuditEntry(userId uuid.UUID, patientId *uuid.UUID, action, entityType string, entityId uuid.UUID, fieldsChanged, ipAddress, userAgent string) *AuditEntry {
	return &AuditEntry{
		Id:            uuid.Must(uuid.NewV7()),
		CreatedAt:     time.Now().UTC(),
		UserId:        userId,
		PatientId:     patientId,
		Action:        action,
		EntityType:    entityType,
		EntityId:      entityId,
		FieldsChanged: fieldsChanged,
		IPAddress:     ipAddress,
		UserAgent:     userAgent,
	}
}
