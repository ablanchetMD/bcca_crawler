package protocols

import (
	"bcca_crawler/api"
	"bcca_crawler/internal/database"
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
)

//Cautions

type CautionReq struct {
	ID          string `json:"id" validate:"omitempty,uuid"`
	Description string `json:"description" validate:"required,min=1,max=500"`
	ProtocolID  string `json:"protocol_id" validate:"omitempty,uuid"`
}

type ChangeProtocolReq struct {
	SelectedProtocolIDs []string `json:"protocol_ids"`
}

type CautionResp struct {
	ID              uuid.UUID             `json:"id"`
	CreatedAt       time.Time             `json:"created_at"`
	UpdatedAt       time.Time             `json:"updated_at"`
	Description     string                `json:"description"`
	LinkedProtocols []api.LinkedProtocols `json:"linked_protocols"`
}

type CautionLike struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Description string
	ProtocolIds json.RawMessage
}

//Eligibility Criteria

type EligibilityCriterionReq struct {
	ID          string `json:"id" validate:"omitempty,uuid"`
	Type        string `json:"type" validate:"required,eligibility_criteria"`
	Description string `json:"description" validate:"required,min=1,max=500"`
	ProtocolID  string `json:"protocol_id" validate:"omitempty,uuid"`
}

type EligibilityCriteriaReq struct {
	EligibilityCriteria []EligibilityCriterionReq `json:"eligibility_criteria"`
}

type EligibilityCriterionResp struct {
	ID              uuid.UUID             `json:"id"`
	Type            string                `json:"type"`
	Description     string                `json:"description"`
	CreatedAt       time.Time             `json:"created_at"`
	UpdatedAt       time.Time             `json:"updated_at"`
	LinkedProtocols []api.LinkedProtocols `json:"linked_protocols"`
}

type EligibilityLike struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Type        database.EligibilityEnum
	Description string
	ProtocolIds interface{}
}

func (e *EligibilityCriterionReq) ToTypeEnum() database.EligibilityEnum {
	return database.EligibilityEnum(strings.ToLower(e.Type))
}

//Precautions

type PrecautionReq struct {
	ID          string `json:"id" validate:"omitempty,uuid"`
	Title       string `json:"title" validate:"required,min=1,max=250"`
	Description string `json:"description" validate:"required,min=1,max=500"`
	ProtocolID  string `json:"protocol_id" validate:"omitempty,uuid"`
}

type PrecautionResp struct {
	ID              uuid.UUID             `json:"id"`
	Title           string                `json:"title"`
	Description     string                `json:"description"`
	CreatedAt       time.Time             `json:"created_at"`
	UpdatedAt       time.Time             `json:"updated_at"`
	LinkedProtocols []api.LinkedProtocols `json:"linked_protocols"`
}

type PrecautionLike struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Title       string
	Description string
	ProtocolIds json.RawMessage
}

//Medications

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
	ID             string   `json:"id" validate:"omitempty,uuid"`
	Name           string   `json:"name" validate:"required,min=1,max=250"`
	Description    string   `json:"description" validate:"omitempty,min=1,max=500"`
	Category       string   `json:"category" validate:"omitempty,min=1,max=50"`
	AlternateNames []string `json:"alternate_names" validate:"omitempty,min=1,max=500"`
}

type PrescriptionReq struct {
	ID           string `json:"id" validate:"omitempty,uuid"`
	MedicationID string `json:"medication_id" validate:"required,uuid"`
	Dose         string `json:"dose" validate:"required"`
	Route        string `json:"route" validate:"required,prescription_route"`
	Frequency    string `json:"frequency" validate:"required"`
	Duration     string `json:"duration" validate:"omitempty"`
	Instructions string `json:"instructions" validate:"omitempty,min=1,max=1000"`
	Renewals     int32  `json:"renewals" validate:"omitempty,min=0,max=50"`
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

type ProtocolMedications struct {
	Premed  []PrescriptionWithMedName `json:"premed"`
	Support []PrescriptionWithMedName `json:"support"`
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

//labs

type LabResp struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Description  string    `json:"description"`
	FormUrl      string    `json:"form_url"`
	Unit         string    `json:"unit"`
	LowerLimit   float64   `json:"lower_limit"`
	UpperLimit   float64   `json:"upper_limit"`
	TestCategory string    `json:"test_category"`
}

type ProtocolTestGroup struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Category  string    `json:"category"`
	Comments  string    `json:"comments"`
	Position  int32     `json:"position"`
	Tests     []LabResp `json:"tests"`
}

type LabsByProtocol struct {
	Baseline map[string][]LabSummary `json:"baseline"`
	FollowUp map[string][]LabSummary `json:"followup"`
}

type LabSummary struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

//toxicities

type ToxicityReq struct {
	ID          string             `json:"id" validate:"omitempty,uuid"`
	Title       string             `json:"title" validate:"required"`
	Category    string             `json:"category" validate:"required"`
	Description string             `json:"description" validate:"omitempty,min=1,max=1000"`
	Grades      []ToxicityGradeReq `json:"grades" validate:"required"`
}

type ToxicityGradeReq struct {
	ID          string `json:"id" validate:"omitempty,uuid"`
	Grade       string `json:"grade" validate:"required,grade"`
	Description string `json:"description" validate:"min=1,max=1000"`
}

type ToxicityGradeWithAdjustmentReq struct {
	ID         string `json:"id" validate:"omitempty,uuid"`
	GradeID    string `json:"toxicity_grade_id" validate:"required,uuid"`
	Adjustment string `json:"adjustment" validate:"required,omitempty"`
}

type ToxModificationsReq struct {
	Grades []ToxicityGradeWithAdjustmentReq `json:"grades" validate:"required"`
}

type ToxModReq struct {
	ID              string `json:"id" validate:"omitempty,uuid"`
	ToxicityGradeID string `json:"toxicity_grade_id" validate:"required,uuid"`
	Adjustment      string `json:"adjustment" validate:"required"`
}

// Treatments

type TreatmentReq struct {
	ID                  string `json:"id" validate:"omitempty,uuid"`
	MedicationID        string `json:"medication_id" validate:"required,uuid"`
	Dose                string `json:"dose" validate:"required"`
	Route               string `json:"route" validate:"required,prescription_route"`
	Frequency           string `json:"frequency" validate:"required"`
	Duration            string `json:"duration" validate:"required"`
	AdministrationGuide string `json:"administration_guide" validate:"omitempty,min=1,max=1000"`
}

type CycleReq struct {
	ID            string `json:"id" validate:"omitempty,uuid"`
	Cycle         string `json:"cycle" validate:"required"`
	CycleDuration string `json:"cycle_duration" validate:"omitempty"`
}

type Adjustment struct {
	Adjustment string `json:"adjustment"`
}

type Subcategory struct {
	Subcategory string       `json:"subcategory"`
	Adjustments []Adjustment `json:"adjustments"`
}

type Category struct {
	Category      string        `json:"category"`
	Subcategories []Subcategory `json:"subcategories"`
}

type MedicationWithModifications struct {
	MedicationID   uuid.UUID  `json:"medication_id"`
	MedicationName string     `json:"medication_name"`
	Categories     []Category `json:"categories"`
}
