package api

import (
	"bcca_crawler/internal/database"
	"github.com/google/uuid"
	"encoding/json"
	"time"
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
	Id      uuid.UUID   `json:"id"`
	Title   string `json:"Title"`
	Authors string `json:"Authors"`
	Journal string `json:"Journal"`
	Year    string `json:"Year"`
	Pmid    string `json:"Pmid"`
	Doi     string `json:"Doi"`
}

type LinkedProtocols struct {
	ID string `json:"id"`
	Code string `json:"code"`
}

type Physician struct {
	Id 	  uuid.UUID `json:"id"`
	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
}

type SummaryProtocol struct {
	Id		 uuid.UUID `json:"id"`
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
	Id 			uuid.UUID `json:"id"`
	Type        database.EligibilityEnum `json:"Type"`
	Description string `json:"Description"`
}

type ProtocolPrecaution struct {
	Id 		uuid.UUID `json:"id"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
}

type ProtocolCaution struct {	
	Id 		uuid.UUID `json:"id"`
	Description string `json:"Description"`
}

type Treatment struct {
	Id                    uuid.UUID              `json:"id"`
	MedicationName        string             	 `json:"Medication"`
	MedicationId		  uuid.UUID              `json:"MedicationID"`
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
	Id          uuid.UUID 		`json:"id"`
	Description string 			`json:"Description"`
	Adjustment  string			`json:"Adjustment"`
}


type ToxicityModification struct {
	Id          uuid.UUID `json:"id"`
	GradeId     uuid.UUID `json:"GradeId"`	
	Grade       string    `json:"Grade"`
	GradeDescription string `json:"GradeDescription"`
	Adjustment  string `json:"Adjustment"`
}

type Toxicity struct {	
	Id            uuid.UUID `json:"id"`
	Title         string `json:"Title"`
	Description   string `json:"Description"`
	Category      string `json:"Category"`
	Modifications []ToxicityModification `json:"Modifications"`
}

type ToxicityGrade struct {
    ID          uuid.UUID  `json:"id"`
    CreatedAt   time.Time  `json:"created_at"`
    UpdatedAt   time.Time  `json:"updated_at"`
    Grade       string     `json:"grade"`
    Description string     `json:"description"`
}

type ToxicityWithGrades struct {
    ID          uuid.UUID        `json:"id"`
    CreatedAt   time.Time        `json:"created_at"`
    UpdatedAt   time.Time        `json:"updated_at"`
    Title       string           `json:"title"`
    Category    string           `json:"category"`
    Description string           `json:"description"`
    Grades      []ToxicityGrade  `json:"grades"`
}

type ToxicityGradeWithAdjustment struct {
    ID          uuid.UUID  `json:"id"`
    CreatedAt   time.Time  `json:"createdAt"`
    UpdatedAt   time.Time  `json:"updatedAt"`
    Grade       string     `json:"grade"`
    Description string     `json:"description"`
    Adjustment  *string    `json:"adjustment"` // pointer because it might be null
}

type ToxicityWithGradesAndAdjustments struct {
    ID          uuid.UUID                     `json:"id"`
    CreatedAt   time.Time                     `json:"createdAt"`
    UpdatedAt   time.Time                     `json:"updatedAt"`
    Title       string                        `json:"title"`
    Category    string                        `json:"category"`
    Description string                        `json:"description"`
    Grades      []ToxicityGradeWithAdjustment `json:"grades"`
}


// Conversion function
func MapToToxicityWithGrades(row database.GetToxicitiesWithGradesRow) (ToxicityWithGrades, error) {
    var grades []ToxicityGrade
    if err := json.Unmarshal(row.Grades, &grades); err != nil {
        return ToxicityWithGrades{}, err
    }
    
    return ToxicityWithGrades{
        ID:          row.ID,
        CreatedAt:   row.CreatedAt,
        UpdatedAt:   row.UpdatedAt,
        Title:       row.Title,
        Category:    row.Category,
        Description: row.Description,
        Grades:      grades,
    }, nil
}

func MapToToxicityWithGradesOne(row database.GetToxicityByIDRow) (ToxicityWithGrades, error) {
    var grades []ToxicityGrade
    if err := json.Unmarshal(row.Grades, &grades); err != nil {
        return ToxicityWithGrades{}, err
    }
    
    return ToxicityWithGrades{
        ID:          row.ID,
        CreatedAt:   row.CreatedAt,
        UpdatedAt:   row.UpdatedAt,
        Title:       row.Title,
        Category:    row.Category,
        Description: row.Description,
        Grades:      grades,
    }, nil
}

func MapToToxicityWithGradesAndAdjustments(row database.GetToxicitiesWithGradesAndAdjustmentsRow) (ToxicityWithGradesAndAdjustments, error) {
    var grades []ToxicityGradeWithAdjustment
    if err := json.Unmarshal(row.Grades, &grades); err != nil {
        return ToxicityWithGradesAndAdjustments{}, err
    }
    
    return ToxicityWithGradesAndAdjustments{
        ID:          row.ID,
        CreatedAt:   row.CreatedAt,
        UpdatedAt:   row.UpdatedAt,
        Title:       row.Title,
        Category:    row.Category,
        Description: row.Description,
        Grades:      grades,
    }, nil
}

func MapToToxicityGrade(tox database.ToxicityGrade) ToxicityGrade {
	return ToxicityGrade{
		ID:          tox.ID,
		CreatedAt:   tox.CreatedAt,
		UpdatedAt:   tox.UpdatedAt,		
		Description: tox.Description,		
		Grade:      string(tox.Grade),
	}
}

type ProtocolCycle struct {
	Id            uuid.UUID `json:"id"`
	Cycle         string         `json:"Cycle"`
	CycleDuration string         `json:"CycleDuration"`
	Treatments    []Treatment    `json:"Treatments"`
}

type MedicationResp struct {
	ID 			string `json:"id"`
	Name 		string `json:"name"`
	Description string `json:"description"`
	CreatedAt 	string `json:"created_at"`
	UpdatedAt 	string `json:"updated_at"`
	Category 	string `json:"category"`
	AlternateNames []string `json:"alternate_names"`
}

func MapMedication(src database.Medication) MedicationResp {
	return MedicationResp{
		ID:          src.ID.String(),
		Name:        src.Name,
		CreatedAt:   src.CreatedAt.Format("2006-01-02"),
		UpdatedAt:   src.UpdatedAt.Format("2006-01-02"),
		Description: src.Description,
		Category:    src.Category,
		AlternateNames: src.AlternateNames,
	}
}

type MedModificationResp struct {
	ID 			string `json:"id"`
	MedicationID string `json:"medication_id"`

	Category 	string `json:"category"`
	Subcategory string `json:"subcategory"`
	Adjustment 	string `json:"adjustment"`
}

func MapMedModification(src database.MedicationModification) MedModificationResp {
	return MedModificationResp{
		ID:          src.ID.String(),
		MedicationID: src.MedicationID.String(),
		Category:    src.Category,
		Subcategory: src.Subcategory,
		Adjustment:  src.Adjustment,
	}
}

type PrescriptionResp struct {
	ID 			string `json:"id"`
	MedicationID string `json:"medication_id"`
	MedicationName string `json:"medication_name"`
	CreatedAt 	string `json:"created_at"`
	UpdatedAt 	string `json:"updated_at"`
	Dose 		string `json:"dose"`
	Route 		string `json:"route"`
	Frequency 	string `json:"frequency"`
	Duration 	string `json:"duration"`
	Instructions string `json:"instructions"`
	Renewals 	int32 `json:"renewals"`
}

func MapPrescriptionByProtByCat(src database.GetPrescriptionsByProtocolByCategoryRow) PrescriptionResp {
	return PrescriptionResp{
		ID:            src.MedicationPrescriptionID.String(),
		MedicationID:  src.MedicationID.String(),
		MedicationName: src.Name,
		CreatedAt:   src.CreatedAt.Format("2006-01-02"),
		UpdatedAt:   src.UpdatedAt.Format("2006-01-02"),
		Dose:          src.Dose,
		Route:         string(src.Route),
		Frequency:     src.Frequency,
		Duration:      src.Duration,
		Instructions:  src.Instructions,
		Renewals:      src.Renewals,
	}
}

func MapPrescriptionsByID(src database.GetPrescriptionByIDRow) PrescriptionResp {
	return PrescriptionResp{
		ID:            src.MedicationPrescriptionID.String(),
		MedicationID:  src.MedicationID.String(),
		MedicationName: src.Name,
		CreatedAt:   src.CreatedAt.Format("2006-01-02"),
		UpdatedAt:   src.UpdatedAt.Format("2006-01-02"),
		Dose:          src.Dose,
		Route:         string(src.Route),
		Frequency:     src.Frequency,
		Duration:      src.Duration,
		Instructions:  src.Instructions,
		Renewals:      src.Renewals,
	}
}

func MapPrescription(src database.GetPrescriptionsRow) PrescriptionResp {
	return PrescriptionResp{
		ID:            src.MedicationPrescriptionID.String(),
		MedicationID:  src.MedicationID.String(),
		MedicationName: src.Name,
		CreatedAt:   src.CreatedAt.Format("2006-01-02"),
		UpdatedAt:   src.UpdatedAt.Format("2006-01-02"),
		Dose:          src.Dose,
		Route:         string(src.Route),
		Frequency:     src.Frequency,
		Duration:      src.Duration,
		Instructions:  src.Instructions,
		Renewals:      src.Renewals,
	}
}

func mapArticleRef(src database.ArticleReference) ArticleReference {
	
	return ArticleReference{
		Id:      src.ID,
		Title:        src.Title,
		Authors: src.Authors,
		Journal: src.Journal,
		Year:      src.Year,
		Pmid:    src.Pmid,
		Doi:    src.Doi,
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

type LabResp struct {
	ID 			string `json:"id"`
	Name 		string `json:"name"`
	CreatedAt 	string `json:"created_at"`
	UpdatedAt 	string `json:"updated_at"`
	Description string `json:"description"`
	FormUrl 	string `json:"form_url"`
	Unit 		string `json:"unit"`
	LowerLimit 	float64 `json:"lower_limit"`
	UpperLimit 	float64 `json:"upper_limit"`
	TestCategory string `json:"test_category"`	
}

func MapLab(src database.Test) LabResp {
	return LabResp{
		ID:          src.ID.String(),
		Name:       src.Name,
		CreatedAt:   src.CreatedAt.Format("2006-01-02"),
		UpdatedAt:   src.UpdatedAt.Format("2006-01-02"),
		Description: src.Description,
		FormUrl:     src.FormUrl,
		Unit:        src.Unit,
		LowerLimit:  src.LowerLimit,
		UpperLimit:  src.UpperLimit,
		TestCategory: src.TestCategory,
	}
}

func MapEligibilityCriterion(src database.ProtocolEligibilityCriterium) ProtocolEligibilityCriterion {
	return ProtocolEligibilityCriterion{
		Id: 		src.ID,
		Type:        src.Type,
		Description: src.Description,
	}
}

func MapPrecaution(src database.ProtocolPrecaution) ProtocolPrecaution {
	return ProtocolPrecaution{
		Id: 		src.ID,
		Title:       src.Title,
		Description: src.Description,
	}
}

func MapCaution(src database.ProtocolCaution) ProtocolCaution {
	return ProtocolCaution{
		Id: 		src.ID,
		Description: src.Description,
	}
}

func MapCycle(src database.ProtocolCycle) ProtocolCycle {
	return ProtocolCycle{
		Id:            src.ID,
		Cycle:         src.Cycle,
		CycleDuration: src.CycleDuration,
	}
}


func MapTreatment(src database.ProtocolTreatment) Treatment {
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
			Grade:            string(row.ToxicityGrade),
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



func MapToMedicationModifications(rows []database.GetMedicationModificationsByProtocolRow) []MedicationModification {
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
			Description: row.ModificationSubcategory,
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
