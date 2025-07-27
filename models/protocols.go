package models

import (
	"time"

	"github.com/google/uuid"
)

type LinkedProtocols struct {
	ID        uuid.UUID `json:"id"`
	Code      string    `json:"code"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CycleReq struct {
	ID            uuid.UUID `json:"id" validate:"omitempty,uuid"`
	Cycle         string    `json:"cycle" validate:"required"`
	CycleDuration string    `json:"cycle_duration" validate:"omitempty"`
}
