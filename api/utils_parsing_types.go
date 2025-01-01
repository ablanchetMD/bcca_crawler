package api

import (
	"bcca_crawler/internal/database"
	"github.com/google/uuid"
)

type ProtocolPayload struct {	
	ProtocolSummary             SummaryProtocol                 `json:"SummaryProtocol"`
	ProtocolEligibilityCriteria []ProtocolEligibilityCriterion `json:"ProtocolEligibilityCriteria"`
	ProtocolPrecautions        []ProtocolPrecaution        		`json:"ProtocolPrecautions"`
	ProtocolCautions		   []ProtocolCaution			   `json:"ProtocolCautions"`
	Tests                      Tests                     	  `json:"Tests"`
	ProtocolCycles             []ProtocolCycle             `json:"ProtocolCycles"`	
	Toxicities			       []Toxicity      				`json:"Toxicities"`
	TreatmentModifications	   []MedicationModification    `json:"TreatmentModifications"`
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
	RevisedOn  string `json:"RevisedOn"`
	ActivatedOn string `json:"ActivatedOn"`
	ProtocolUrl string `json:"ProtocolUrl"`
	HandOutUrl string `json:"HandOutUrl"`
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

type MedicationModification struct {
	MedicationId			uuid.UUID 						`json:"MedicationId"`	
	Medication  			string 							`json:"Medication"`
	ModificationCategory 	[]ModificationCategory 			`json:"ModificationCategory"`
}

type ModificationCategory struct {
	Category 		string			 `json:"Category"`	
	Modifications 	[]Modifications	 `json:"Modifications"`
}

type Modifications struct {
	Id          uuid.UUID 		`json:"Id"`
	Description string 			`json:"Description"`
	Adjustment  string			`json:"Adjustment"`
}


type ToxicityModification struct {
	Id          uuid.UUID `json:"Id"`
	GradeId     uuid.UUID `json:"GradeId"`	
	Grade       string    `json:"Grade"`
	GradeDescription string `json:"GradeDescription"`
	Adjustment  string `json:"Adjustment"`
}

type Toxicity struct {	
	Id            uuid.UUID `json:"Id"`
	Title         string `json:"Title"`
	Description   string `json:"Description"`
	Category      string `json:"Category"`
	Modifications []ToxicityModification `json:"Modifications"`
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
		ProtocolUrl: src.ProtocolUrl,
		HandOutUrl: src.PatientHandoutUrl,
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


func mapTest(src []database.Test) []string {
	var tests []string
	for _, t := range src {
		tests = append(tests, t.Name)
	}
	return tests
}


func mapToToxicities(rows []database.GetToxicityModificationByProtocolRow) []Toxicity {
    // Map to group toxicities by their ID
    toxicityMap := make(map[uuid.UUID]*Toxicity)

    for _, row := range rows {
        // Check if the toxicity already exists in the map
        if _, exists := toxicityMap[row.ID]; !exists {
            // Add a new toxicity to the map
            toxicityMap[row.ID] = &Toxicity{
                Id:          row.ID,
                Title:       row.ToxicityTitle,
                Description: row.ToxicityGradeDescription, // Assuming this maps correctly
                Modifications: []ToxicityModification{},
            }
        }

        // Add the modification to the existing toxicity
        toxicityMap[row.ID].Modifications = append(toxicityMap[row.ID].Modifications, ToxicityModification{
            Id:               row.ID,
            GradeId:          row.ToxicityGradeID,
            Grade:            row.ToxicityGrade,
            GradeDescription: row.ToxicityGradeDescription,
            Adjustment:       row.Adjustment,
        })
    }

    // Convert the map to a slice
    toxicities := make([]Toxicity, 0, len(toxicityMap))
    for _, toxicity := range toxicityMap {
        toxicities = append(toxicities, *toxicity)
    }

    return toxicities
}



func mapToMedicationModifications(rows []database.GetMedicationModificationsByProtocolRow) []MedicationModification {
	// Map to group medications by MedicationID
	medicationMap := make(map[uuid.UUID]*MedicationModification)

	for _, row := range rows {
		// Check if the medication already exists in the map
		if _, exists := medicationMap[row.MedicationID]; !exists {
			// Add a new medication to the map
			medicationMap[row.MedicationID] = &MedicationModification{
				MedicationId: row.MedicationID,
				Medication:   row.Name,
				ModificationCategory: []ModificationCategory{},
			}
		}

		// Find or create the ModificationCategory within the Medication
		var category *ModificationCategory
		for i := range medicationMap[row.MedicationID].ModificationCategory {
			if medicationMap[row.MedicationID].ModificationCategory[i].Category == row.ModificationCategory {
				category = &medicationMap[row.MedicationID].ModificationCategory[i]
				break
			}
		}

		// If the category doesn't exist, create it
		if category == nil {
			category = &ModificationCategory{
				Category:     row.ModificationCategory,
				Modifications: []Modifications{},
			}
			medicationMap[row.MedicationID].ModificationCategory = append(medicationMap[row.MedicationID].ModificationCategory, *category)
		}

		// Add the modification to the category
		category.Modifications = append(category.Modifications, Modifications{
			Id:          row.ModificationID,
			Description: row.ModificationDescription,
			Adjustment:  row.Adjustment,
		})
	}

	// Convert the map to a slice
	medications := make([]MedicationModification, 0, len(medicationMap))
	for _, medication := range medicationMap {
		medications = append(medications, *medication)
	}

	return medications
}
