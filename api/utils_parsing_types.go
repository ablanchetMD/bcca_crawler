package api

import (
	"bcca_crawler/internal/database"
	"github.com/google/uuid"
)

type ProtocolPayload struct {	
	ProtocolSummary             SummaryProtocol                 `json:"SummaryProtocol"`
	ProtocolEligibilityCriteria []ProtocolEligibilityCriterion `json:"ProtocolEligibilityCriteria"`
	ProtocolPrecautions        []ProtocolPrecaution        `json:"ProtocolPrecautions"`
	ProtocolCautions		   []ProtocolCaution		   `json:"ProtocolCautions"`
	Tests                      Tests                       `json:"Tests"`
	ProtocolCycles             []ProtocolCycle             `json:"ProtocolCycles"`	
	ToxicityModifications      []ToxicityModification      `json:"ToxicityModifications"`
	Physicians                 []Physician                 `json:"Physicians"`
	ArticleReferences          []ArticleReference          `json:"ArticleReferences"`
}

type ArticleReference struct {
	Id      uuid.UUID   `json:"Id"`
	Title   string `json:"Title"`
	Authors string `json:"Authors"`
	Journal string `json:"Journal"`
	Year    string `json:"Year"`
	Pmid    string `json:"Pmid"`
	Joi     string `json:"Joi"`
}

type Physician struct {
	Id 	  uuid.UUID `json:"Id"`
	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
}

type SummaryProtocol struct {
	Id		 uuid.UUID `json:"Id"`
	TumorGroup string `json:"TumorGroup"`
	Code       string `json:"Code"`
	Name       string `json:"Name"`
	Tags       []string `json:"Tags"`
	Notes      string `json:"Notes"`
}

type ProtocolEligibilityCriterion struct {
	Id 		uuid.UUID `json:"Id"`
	Type        string `json:"Type"`
	Description string `json:"Description"`
}

type ProtocolPrecaution struct {
	Id 		uuid.UUID `json:"Id"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
}

type ProtocolCaution struct {	
	Id 		uuid.UUID `json:"Id"`
	Description string `json:"Description"`
}

type Treatment struct {
	Id                    uuid.UUID              `json:"Id"`
	MedicationName        string             	 `json:"Medication"`
	MedicationId		  uuid.UUID              `json:"MedicationId"`
	Dose                  string                 `json:"Dose"`
	Route                 string                 `json:"Route"`
	Frequency             string                 `json:"Frequency"`
	Duration              string                 `json:"Duration"`
	AdministrationGuide   string                 `json:"AdministrationGuide"`
	TreatmentModifications []TreatmentModification `json:"TreatmentModifications"`
}

type Tests struct {
	Baseline BaselineTests `json:"Baseline"`
	FollowUp FollowUpTests `json:"FollowUp"`
}

type BaselineTests struct {
	RequiredBeforeTreatment []string `json:"RequiredBeforeTreatment"`
	RequiredButCanProceed   []string `json:"RequiredButCanProceed"`
	IfClinicallyIndicated   []string `json:"IfClinicallyIndicated"`
}

type FollowUpTests struct {
	Required               []string `json:"Required"`
	IfClinicallyIndicated  []string `json:"IfClinicallyIndicated"`
}

type TreatmentModification struct {
	Id          uuid.UUID `json:"Id"`
	Category    string `json:"Category"`
	Description string `json:"Description"`
	Adjustement string `json:"Adjustement"`
}

type ToxicityModification struct {
	Id          uuid.UUID `json:"Id"`
	Title       string `json:"Title"`
	Grade       string `json:"Grade"`
	Adjustement string `json:"Adjustement"`
}

type ProtocolCycle struct {
	Id            uuid.UUID `json:"Id"`
	Cycle         string         `json:"Cycle"`
	CycleDuration string         `json:"CycleDuration"`
	Treatments    []Treatment    `json:"Treatments"`
}



func mapArticleRef(src database.ArticleReference) ArticleReference {
	
	return ArticleReference{
		Id:      src.ID,
		Title:        src.Title,
		Authors: src.Authors,
		Journal: src.Journal,
		Year:      src.Year,
		Pmid:    src.Pmid,
		Joi:    src.Joi,
	}
}

func mapPhysician(src database.Physician) Physician {
	return Physician{
		Id: 	  src.ID,
		FirstName: src.FirstName,
		LastName:  src.LastName,
	}
}

func mapSummaryProtocol(src database.Protocol) SummaryProtocol {
	return SummaryProtocol{
		Id:		 src.ID,
		TumorGroup: src.TumorGroup,
		Code:       src.Code,
		Name:       src.Name,
		Tags:       src.Tags,
		Notes:      src.Notes,
	}
}

func mapEligibilityCriterion(src database.ProtocolEligibilityCriterium) ProtocolEligibilityCriterion {
	return ProtocolEligibilityCriterion{
		Id: 		src.ID,
		Type:        src.Type,
		Description: src.Description,
	}
}

func mapPrecaution(src database.ProtocolPrecaution) ProtocolPrecaution {
	return ProtocolPrecaution{
		Id: 		src.ID,
		Title:       src.Title,
		Description: src.Description,
	}
}

func mapCaution(src database.ProtocolCaution) ProtocolCaution {
	return ProtocolCaution{
		Id: 		src.ID,
		Description: src.Description,
	}
}

func mapCycle(src database.ProtocolCycle) ProtocolCycle {
	return ProtocolCycle{
		Id:            src.ID,
		Cycle:         src.Cycle,
		CycleDuration: src.CycleDuration,
	}
}


func mapTreatment(src database.ProtocolTreatment) Treatment {
	return Treatment{
		Id:                    src.ID,
		MedicationId: 		   src.Medication,		
		Dose:                  src.Dose,
		Route:                 src.Route,
		Frequency:             src.Frequency,
		Duration:              src.Duration,
		AdministrationGuide:   src.AdministrationGuide,
	}
}

func mapTreatmentModification(src database.TreatmentModification) TreatmentModification {
	return TreatmentModification{
		Id:          src.ID,
		Category:    src.Category,
		Description: src.Description,
		Adjustement: src.Adjustement,
	}
}

func mapTest(src []database.Test) []string {
	var tests []string
	for _, t := range src {
		tests = append(tests, t.Name)
	}
	return tests
}

func mapToxicityModification(src database.ToxicityModification) ToxicityModification {
	return ToxicityModification{
		Id:          src.ID,
		Title:       src.Title,
		Grade:       src.Grade,
		Adjustement: src.Adjustement,
	}
}
