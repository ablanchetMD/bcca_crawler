package models

import (
	"encoding/json"
	"time"
	"github.com/google/uuid"
)

type PrecautionReq struct {
	ID          uuid.UUID  `json:"id" validate:"omitempty,uuid"`
	Title       string `json:"title" validate:"required,min=1,max=250"`
	Description string `json:"description" validate:"required,min=1,max=500"`
	ProtocolID  uuid.UUID  `json:"protocol_id" validate:"omitempty,uuid"`
}

type PrecautionResp struct {
	ID              uuid.UUID             `json:"id"`
	Title           string                `json:"title"`
	Description     string                `json:"description"`
	CreatedAt       time.Time             `json:"created_at"`
	UpdatedAt       time.Time             `json:"updated_at"`
	LinkedProtocols []LinkedProtocols `json:"linked_protocols"`
}

type PrecautionLike struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Title       string
	Description string
	ProtocolIds json.RawMessage
}