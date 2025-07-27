package models

import (
	"bcca_crawler/internal/database"
	"strings"
	"time"

	"github.com/google/uuid"
)

type EligibilityCriterionReq struct {
	ID          uuid.UUID `json:"id" validate:"omitempty,uuid"`
	Type        string    `json:"type" validate:"required,eligibility_criteria"`
	Description string    `json:"description" validate:"required,min=1,max=500"`
	ProtocolID  string    `json:"protocol_id" validate:"omitempty,uuid"`
}

type EligibilityCriteriaReq struct {
	EligibilityCriteria []EligibilityCriterionReq `json:"eligibility_criteria"`
}

type EligibilityCriterionResp struct {
	ID              uuid.UUID         `json:"id"`
	Type            string            `json:"type"`
	Description     string            `json:"description"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
	LinkedProtocols []LinkedProtocols `json:"linked_protocols"`
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
