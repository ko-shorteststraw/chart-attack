package interfaces

import (
	"context"

	"github.com/kendall/chart-attack/internal/application/command"
	"github.com/kendall/chart-attack/internal/application/query"
)

type PatientService interface {
	AdmitPatient(ctx context.Context, cmd command.AdmitPatientCommand) (*command.AdmitPatientResult, error)
	GetPatientById(ctx context.Context, q query.GetPatientByIdQuery) (*query.GetPatientByIdResult, error)
	GetPatientCensus(ctx context.Context, q query.GetPatientCensusQuery) (*query.GetPatientCensusResult, error)
	UpdateDiagnoses(ctx context.Context, cmd command.UpdateDiagnosesCommand) error
	AddEmergencyContact(ctx context.Context, cmd command.AddEmergencyContactCommand) (*command.AddEmergencyContactResult, error)
	RemoveEmergencyContact(ctx context.Context, cmd command.RemoveEmergencyContactCommand) error
	GetPatientEmergencyContacts(ctx context.Context, q query.GetPatientEmergencyContactsQuery) (*query.GetPatientEmergencyContactsResult, error)
}
