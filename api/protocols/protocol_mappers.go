package protocols

import (
	"bcca_crawler/internal/database"
	"bcca_crawler/api"
)

//Cautions

func mapToCautionResp[T any](row T) CautionLike {
	switch r := any(row).(type) {
	case database.GetCautionWithProtocolsRow:
		return CautionLike(r)
	case database.GetCautionByIDWithProtocolsRow:
		return CautionLike(r)
	default:
		panic("unsupported row type")
	}
}

func MapCautionWithProtocols[T any](r T) (CautionResp,error) {
	src := mapToCautionResp(r)
	linkedProtocols, err := api.ParseLinkedProtocols(src.ProtocolIds)
		if err != nil {			
			return CautionResp{},err
		}

	return_item := CautionResp{
		ID:          src.ID,
		CreatedAt:	src.CreatedAt,
		UpdatedAt: 	src.UpdatedAt,
		Description: src.Description,
		LinkedProtocols: linkedProtocols,
	}

	return return_item,nil
}

//Eligibility

func mapToEligibilityResp[T any](row T) EligibilityLike {
	switch r := any(row).(type) {
	case database.GetElibilityCriteriaRow:
		return EligibilityLike(r)
	case database.GetEligibilityCriteriaByIDRow:
		return EligibilityLike(r)
	case database.GetEligibilityCriteriaByTypeRow:
		return EligibilityLike(r)
	default:
		panic("unsupported row type")
	}
}

func MapEligibilityWithProtocols[T any](r T) (EligibilityCriterionResp,error) {
	src := mapToEligibilityResp(r)
	linkedProtocols, err := api.ParseLinkedProtocols(src.ProtocolIds)
		if err != nil {			
			return EligibilityCriterionResp{},err
		}

	return_item := EligibilityCriterionResp{
		ID:          src.ID,
		CreatedAt:	src.CreatedAt,
		UpdatedAt: 	src.UpdatedAt,
		Description: src.Description,
		Type:		string(src.Type),
		LinkedProtocols: linkedProtocols,
	}

	return return_item,nil
}

//Precaution

func mapToPrecautionResponse[T any](row T) PrecautionLike {
	switch r := any(row).(type) {
	case database.GetPrecautionWithProtocolsRow:
		return PrecautionLike(r)
	case database.GetPrecautionByIDWithProtocolsRow:
		return PrecautionLike(r)	
	default:
		panic("unsupported row type")
	}
}

func MapPrecautionWithProtocols[T any](r T) (PrecautionResp,error) {
	src := mapToPrecautionResponse(r)
	linkedProtocols, err := api.ParseLinkedProtocols(src.ProtocolIds)
		if err != nil {			
			return PrecautionResp{},err
		}

	return_item := PrecautionResp{
		ID:          src.ID,
		CreatedAt:	src.CreatedAt,
		UpdatedAt: 	src.UpdatedAt,
		Description: src.Description,
		Title:		src.Title,
		LinkedProtocols: linkedProtocols,
	}

	return return_item,nil
}

//Medications

func MapMedication(src database.Medication) MedicationResp {
	return MedicationResp{
		ID:          src.ID,
		Name:        src.Name,
		CreatedAt:   src.CreatedAt,
		UpdatedAt:   src.UpdatedAt,
		Description: src.Description,
		Category:    src.Category,
		AlternateNames: src.AlternateNames,
	}
}

func mapToPrescriptionResp[T any](row T) PrescriptionLike {
	switch r := any(row).(type) {
	case database.GetPrescriptionByIDRow:
		return PrescriptionLike(r)
	case database.GetPrescriptionsByMedRow:
		return PrescriptionLike(r)
	case database.GetPrescriptionsRow:
		return PrescriptionLike(r)
	case database.GetPrescriptionsByProtocolByCategoryRow:
		return PrescriptionLike(r)
	default:
		panic("unsupported row type")
	}
}

func MapPrescription[T any](r T) PrescriptionResp {
	src := mapToPrescriptionResp(r)
	return PrescriptionResp{
		ID:            	src.MedicationPrescriptionID,
		MedicationID:  	src.MedicationID,
		MedicationName: src.Name,
		CreatedAt:   	src.CreatedAt,
		UpdatedAt:   	src.UpdatedAt,
		Dose:          	src.Dose,
		Route:         	string(src.Route),
		Frequency:     	src.Frequency,
		Duration:      	src.Duration,
		Instructions:  	src.Instructions,
		Renewals:     	 src.Renewals,
	}
}


//Labs

func MapLab(src database.Test) LabResp {
	return LabResp{
		ID:          src.ID,
		Name:       src.Name,
		CreatedAt:   src.CreatedAt,
		UpdatedAt:   src.UpdatedAt,
		Description: src.Description,
		FormUrl:     src.FormUrl,
		Unit:        src.Unit,
		LowerLimit:  src.LowerLimit,
		UpperLimit:  src.UpperLimit,
		TestCategory: src.TestCategory,
	}
}