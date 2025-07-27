package models

import (
	"time"
	"github.com/google/uuid"
)

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

type LabSummary struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}
