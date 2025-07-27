package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type CautionReq struct {
	ID          uuid.UUID `json:"id" validate:"omitempty,uuid"`
	Description string    `json:"description" validate:"required,min=1,max=500"`
	ProtocolID  string    `json:"protocol_id" validate:"omitempty,uuid"`
}

type ChangeProtocolReq struct {
	SelectedProtocolIDs []string `json:"protocol_ids"`
}

type CautionResp struct {
	ID              uuid.UUID         `json:"id"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
	Description     string            `json:"description"`
	LinkedProtocols []LinkedProtocols `json:"linked_protocols"`
}

type CautionLike struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Description string
	ProtocolIds json.RawMessage
}
