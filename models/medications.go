package models

import (
	"bcca_crawler/internal/database"
	"strings"
	"time"

	"github.com/google/uuid"
)

type MedicationResp struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Category       string    `json:"category"`
	AlternateNames []string  `json:"alternate_names"`
}

type MedReq struct {
	ID             uuid.UUID `json:"id" validate:"omitempty,uuid"`
	Name           string    `json:"name" validate:"required,min=1,max=250"`
	Description    string    `json:"description" validate:"omitempty,min=1,max=500"`
	Category       string    `json:"category" validate:"omitempty,min=1,max=50"`
	AlternateNames []string  `json:"alternate_names" validate:"omitempty,min=1,max=500"`
}

type PrescriptionReq struct {
	ID           uuid.UUID `json:"id" validate:"omitempty,uuid"`
	MedicationID string    `json:"medication_id" validate:"required,uuid"`
	Dose         string    `json:"dose" validate:"required"`
	Route        string    `json:"route" validate:"required,prescription_route"`
	Frequency    string    `json:"frequency" validate:"required"`
	Duration     string    `json:"duration" validate:"omitempty"`
	Instructions string    `json:"instructions" validate:"omitempty,min=1,max=1000"`
	Renewals     int32     `json:"renewals" validate:"omitempty,min=0,max=50"`
}

func (e *PrescriptionReq) ToRouteEnum() database.PrescriptionRouteEnum {
	return database.PrescriptionRouteEnum(strings.ToLower(e.Route))
}

func (e *TreatmentReq) ToRouteEnum() database.PrescriptionRouteEnum {
	return database.PrescriptionRouteEnum(strings.ToLower(e.Route))
}

type PrescriptionLike struct {
	MedicationID             uuid.UUID
	Name                     string
	Description              string
	Category                 string
	AlternateNames           []string
	MedicationPrescriptionID uuid.UUID
	Dose                     string
	CreatedAt                time.Time
	UpdatedAt                time.Time
	Route                    database.PrescriptionRouteEnum
	Frequency                string
	Duration                 string
	Instructions             string
	Renewals                 int32
}

type PrescriptionResp struct {
	ID                    uuid.UUID `json:"id"`
	MedicationID          uuid.UUID `json:"medication_id"`
	MedicationName        string    `json:"medication_name"`
	MedicationDescription string    `json:"medication_description"`
	MedicationCategory    string    `json:"medication_category"`
	MedicationAlternates  []string  `json:"medication_alternate_names"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
	Dose                  string    `json:"dose"`
	Route                 string    `json:"route"`
	Frequency             string    `json:"frequency"`
	Duration              string    `json:"duration"`
	Instructions          string    `json:"instructions"`
	Renewals              int32     `json:"renewals"`
}

type PrescriptionWithMedName struct {
	ID                    uuid.UUID `json:"id"`
	Dose                  string    `json:"dose"`
	Route                 string    `json:"route"`
	Frequency             string    `json:"frequency"`
	Duration              string    `json:"duration"`
	Instructions          string    `json:"instructions"`
	Renewals              int       `json:"renewals"`
	MedicationID          uuid.UUID `json:"medication_id"`
	MedicationName        string    `json:"medication_name"`
	MedicationDescription string    `json:"medication_description"`
	MedicationAlternates  []string  `json:"medication_alternate_names"`
	MedicationCategory    string    `json:"medication_category"`
}

type ProtocolMedGroup struct {
	ID          uuid.UUID                 `json:"id"`
	CreatedAt   time.Time                 `json:"created_at"`
	UpdatedAt   time.Time                 `json:"updated_at"`
	Category    string                    `json:"category"`
	Comments    string                    `json:"comments"`
	Medications []PrescriptionWithMedName `json:"medications"`
}

type TreatmentReq struct {
	ID                  uuid.UUID `json:"id" validate:"omitempty,uuid"`
	MedicationID        uuid.UUID `json:"medication_id" validate:"required,uuid"`
	Dose                string    `json:"dose" validate:"required"`
	Route               string    `json:"route" validate:"required,prescription_route"`
	Frequency           string    `json:"frequency" validate:"required"`
	Duration            string    `json:"duration" validate:"required"`
	AdministrationGuide string    `json:"administration_guide" validate:"omitempty,min=1,max=1000"`
}
