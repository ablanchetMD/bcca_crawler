package models

import "github.com/google/uuid"

//toxicities

type ToxicityReq struct {
	ID          uuid.UUID             `json:"id" validate:"omitempty,uuid"`
	Title       string             `json:"title" validate:"required"`
	Category    string             `json:"category" validate:"required"`
	Description string             `json:"description" validate:"omitempty,min=1,max=1000"`
	Grades      []ToxicityGradeReq `json:"grades" validate:"required"`
}

type ToxicityGradeReq struct {
	ID          uuid.UUID  `json:"id" validate:"omitempty,uuid"`
	Grade       string `json:"grade" validate:"required,grade"`
	Description string `json:"description" validate:"min=1,max=1000"`
}

type ToxicityGradeWithAdjustmentReq struct {
	ID         uuid.UUID  `json:"id" validate:"omitempty,uuid"`
	GradeID    string `json:"toxicity_grade_id" validate:"required,uuid"`
	Adjustment string `json:"adjustment" validate:"required,omitempty"`
}

type ToxModificationsReq struct {
	Grades []ToxicityGradeWithAdjustmentReq `json:"grades" validate:"required"`
}

type ToxModReq struct {
	ID              uuid.UUID  `json:"id" validate:"omitempty,uuid"`
	ToxicityGradeID string `json:"toxicity_grade_id" validate:"required,uuid"`
	Adjustment      string `json:"adjustment" validate:"required"`
}