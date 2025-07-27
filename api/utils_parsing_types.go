package api

import (
	"bcca_crawler/internal/database"
	"bcca_crawler/models"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ProtocolPayload struct {
	ProtocolSummary             SummaryProtocol                `json:"SummaryProtocol"`
	ProtocolEligibilityCriteria []ProtocolEligibilityCriterion `json:"ProtocolEligibilityCriteria"`
	ProtocolPrecautions         []ProtocolPrecaution           `json:"ProtocolPrecautions"`
	ProtocolCautions            []ProtocolCaution              `json:"ProtocolCautions"`
	Tests                       LabsByProtocol                 `json:"Tests"`
	ProtocolCycles              []ProtocolCycle                `json:"ProtocolCycles"`
	Toxicities                  []Toxicity                     `json:"Toxicities"`
	TreatmentModifications      []MedicationModification       `json:"TreatmentModifications"`
	Physicians                  []Physician                    `json:"Physicians"`
	ArticleReferences           []ArticleReference             `json:"ArticleReferences"`
}

type ProtocolSumPayload struct {
	ProtocolSummary             SummaryProtocol                    `json:"protocol_summary"`
	ProtocolEligibilityCriteria []ProtocolEligibilityCriterion     `json:"eligibility_criteria"`
	ProtocolPrecautions         []ProtocolPrecaution               `json:"precautions"`
	ProtocolCautions            []ProtocolCaution                  `json:"cautions"`
	Tests                       []models.ProtocolTestGroup         `json:"test_groups"`
	ProtocolMeds                []models.ProtocolMedGroup          `json:"medication_groups"`
	ProtocolCycles              []ProtocolCycle                    `json:"cycles"`
	Toxicities                  []ToxicityWithGradesAndAdjustments `json:"toxicities"`
	TreatmentModifications      []MedicationWithModifications      `json:"treatment_modifications"`
	Physicians                  []Physician                        `json:"physicians"`
	ArticleReferences           []ArticleReference                 `json:"article_references"`
}

type ArticleReference struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Authors   string    `json:"authors"`
	Journal   string    `json:"journal"`
	Year      string    `json:"year"`
	Pmid      string    `json:"pmid"`
	Doi       string    `json:"doi"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProtocolMedications struct {
	Premed  []PrescriptionWithMedName `json:"premed"`
	Support []PrescriptionWithMedName `json:"support"`
}

type PrescriptionWithMedName struct {
	ID                    uuid.UUID `json:"id"`
	Dose                  string    `json:"dose"`
	Route                 string    `json:"route"`
	Frequency             string    `json:"frequency"`
	Duration              string    `json:"duration"`
	Instructions          string    `json:"instructions"`
	Renewals              int       `json:"renewals"`
	MedicationID          uuid.UUID `json:"medication_id"`
	MedicationName        string    `json:"medication_name"`
	MedicationDescription string    `json:"medication_description"`
	MedicationCategory    string    `json:"medication_category"`
	MedicationAlternates  []string  `json:"medication_alternate_names"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

type LinkedProtocols struct {
	ID        string    `json:"id"`
	Code      string    `json:"code"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Physician struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SummaryProtocol struct {
	ID          uuid.UUID `json:"id"`
	TumorGroup  string    `json:"tumor_group"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Tags        []string  `json:"tags"`
	Notes       string    `json:"notes"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	RevisedOn   string    `json:"revised_on"`
	ActivatedOn string    `json:"activated_on"`
	ProtocolUrl string    `json:"protocol_url"`
	HandOutUrl  string    `json:"handout_url"`
}

type ProtocolEligibilityCriterion struct {
	ID          uuid.UUID                `json:"id"`
	Type        database.EligibilityEnum `json:"type"`
	Description string                   `json:"description"`
	CreatedAt   time.Time                `json:"created_at"`
	UpdatedAt   time.Time                `json:"updated_at"`
}

type ProtocolPrecaution struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProtocolCaution struct {
	ID          uuid.UUID `json:"id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type LabsByProtocol struct {
	Baseline map[string][]LabSummary `json:"baseline"`
	FollowUp map[string][]LabSummary `json:"followup"`
}

type LabSummary struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

type MedicationModification struct {
	MedicationID         uuid.UUID              `json:"medication_id"`
	Medication           string                 `json:"medication"`
	CreatedAt            time.Time              `json:"created_at"`
	UpdatedAt            time.Time              `json:"updated_at"`
	ModificationCategory []ModificationCategory `json:"modification_category"`
}

type ModificationCategory struct {
	Category      string          `json:"category"`
	Modifications []Modifications `json:"modifications"`
}

type Modifications struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Description string    `json:"description"`
	Adjustment  string    `json:"adjustment"`
}

type ToxicityModification struct {
	ID               uuid.UUID `json:"id"`
	GradeID          uuid.UUID `json:"grade_id"`
	Grade            string    `json:"grade"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	GradeDescription string    `json:"grade_description"`
	Adjustment       string    `json:"adjustment"`
}

type Toxicity struct {
	ID            uuid.UUID              `json:"id"`
	Title         string                 `json:"title"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
	Description   string                 `json:"description"`
	Category      string                 `json:"category"`
	Modifications []ToxicityModification `json:"modifications"`
}

type ToxicityGrade struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Grade       string    `json:"grade"`
	Description string    `json:"description"`
}

type ToxicityWithGrades struct {
	ID          uuid.UUID       `json:"id"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	Title       string          `json:"title"`
	Category    string          `json:"category"`
	Description string          `json:"description"`
	Grades      []ToxicityGrade `json:"grades"`
}

type ToxicityGradeWithAdjustment struct {
	ID           uuid.UUID  `json:"id"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	Grade        string     `json:"grade"`
	Description  string     `json:"description"`
	AdjustmentID *uuid.UUID `json:"adjustment_id"` // pointer because it might be null
	Adjustment   *string    `json:"adjustment"`    // pointer because it might be null
}

type ToxicityWithGradesAndAdjustments struct {
	ID          uuid.UUID                     `json:"id"`
	CreatedAt   time.Time                     `json:"created_at"`
	UpdatedAt   time.Time                     `json:"updated_at"`
	Title       string                        `json:"title"`
	Category    string                        `json:"category"`
	Description string                        `json:"description"`
	Grades      []ToxicityGradeWithAdjustment `json:"grades"`
}

type Subcategory struct {
	Subcategory string `json:"subcategory"`
	Adjustment  string `json:"adjustment"`
}

type Category struct {
	Category      string        `json:"category"`
	Subcategories []Subcategory `json:"subcategories"`
}

type MedicationWithModifications struct {
	MedicationID   uuid.UUID  `json:"medication_id"`
	MedicationName string     `json:"medication_name"`
	Categories     []Category `json:"categories"`
}

type Treatment struct {
	MedicationID          uuid.UUID                      `json:"medication_id"`
	MedicationName        string                         `json:"medication_name"`
	MedicationDescription string                         `json:"medication_description"`
	MedicationCategory    string                         `json:"medication_category"`
	MedicationAlternates  []string                       `json:"medication_alternate_names"`
	ID                    uuid.UUID                      `json:"id"`
	Dose                  string                         `json:"dose"`
	CreatedAt             time.Time                      `json:"created_at"`
	UpdatedAt             time.Time                      `json:"updated_at"`
	Route                 database.PrescriptionRouteEnum `json:"route"`
	Frequency             string                         `json:"frequency"`
	Duration              string                         `json:"duration"`
	AdministrationGuide   string                         `json:"administration_guide"`
}

type Prescription struct {
	MedicationID          uuid.UUID                      `json:"medication_id"`
	MedicationName        string                         `json:"medication_name"`
	MedicationDescription string                         `json:"medication_description"`
	MedicationCategory    string                         `json:"medication_category"`
	MedicationAlternates  []string                       `json:"medication_alternate_names"`
	ID                    uuid.UUID                      `json:"id"`
	Dose                  string                         `json:"dose"`
	CreatedAt             time.Time                      `json:"created_at"`
	UpdatedAt             time.Time                      `json:"updated_at"`
	Route                 database.PrescriptionRouteEnum `json:"route"`
	Frequency             string                         `json:"frequency"`
	Duration              string                         `json:"duration"`
	Instructions          string                         `json:"instructions"`
	Renewals              int32                          `json:"renewals"`
}

type ProtocolCycle struct {
	ID            uuid.UUID   `json:"id"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
	Cycle         string      `json:"cycle"`
	CycleDuration string      `json:"cycle_duration"`
	Treatments    []Treatment `json:"treatments"`
}

type TestGroup struct {
	ID        uuid.UUID    `json:"id"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	Category  string       `json:"category"`
	Comments  string       `json:"comments"`
	Position  int32        `json:"position"`
	Tests     []LabSummary `json:"tests"`
}

type PrescriptionGroup struct {
	ID            uuid.UUID      `json:"id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	Category      string         `json:"category"`
	Comments      string         `json:"comments"`
	Prescriptions []Prescription `json:"prescriptions"`
}

func ParsePostGRESData[T any](data any) ([]T, error) {
	bytes, err := ToJSONBytes(data)
	if err != nil {
		return nil, fmt.Errorf("expected []byte but got %T", data)
	}
	var results []T
	if err := json.Unmarshal(bytes, &results); err != nil {
		return nil, fmt.Errorf("failed to unmarshal into []%T: %w", results, err)
	}
	return results, nil
}

func ToJSONBytes(data any) ([]byte, error) {
	switch v := data.(type) {
	case json.RawMessage:
		return v, nil
	case []byte:
		return v, nil
	case string:
		return []byte(v), nil
	default:
		return json.Marshal(v) // fallback: marshal any other Go value
	}
}

func ToResponseData[T any](data any) (T, error) {
	var response T
	jsonBytes, err := ToJSONBytes(data)
	if err != nil {
		return response, fmt.Errorf("expected []byte, got %T", data)
	}

	if err := json.Unmarshal(jsonBytes, &response); err != nil {
		return response, fmt.Errorf("error unmarshalling %v with error: %v", data, err)
	}

	return response, nil
}

// Conversion function

type ToxicityWithGradesLike struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Title       string
	Category    string
	Description string
	Grades      json.RawMessage
}

func mapToTreatmentLike[T any](row T) Treatment {
	switch r := any(row).(type) {
	case database.GetTreatmentsRow:
		return Treatment(r)
	case database.GetTreatmentsByCycleRow:
		return Treatment(r)
	case database.GetProtocolTreatmentByIDRow:
		return Treatment(r)
	default:
		panic("unsupported row type")
	}
}

func MapTreatment[T any](r T) Treatment {
	src := mapToTreatmentLike(r)
	return Treatment{
		ID:                    src.ID,
		MedicationID:          src.MedicationID,
		MedicationName:        src.MedicationName,
		MedicationDescription: src.MedicationDescription,
		MedicationCategory:    src.MedicationCategory,
		MedicationAlternates:  src.MedicationAlternates,
		CreatedAt:             src.CreatedAt,
		UpdatedAt:             src.UpdatedAt,
		Dose:                  src.Dose,
		Route:                 src.Route,
		Frequency:             src.Frequency,
		Duration:              src.Duration,
		AdministrationGuide:   src.AdministrationGuide,
	}
}

func mapToToxicityWithGrades[T any](row T) ToxicityWithGradesLike {
	switch r := any(row).(type) {
	case database.GetToxicitiesWithGradesRow:
		return ToxicityWithGradesLike(r)
	case database.GetToxicityByIDRow:
		return ToxicityWithGradesLike(r)
	case database.GetToxicitiesWithGradesAndAdjustmentsRow:
		return ToxicityWithGradesLike(r)
	case database.GetToxicitiesWithGradesAndAdjustmentsByProtocolRow:
		return ToxicityWithGradesLike(r)
	default:
		panic("unsupported row type")
	}
}

func MapMedModification2(src database.GetProtocolMedicationsWithModificationsRow) (MedicationWithModifications, error) {
	cats, err := ParsePostGRESData[Category](src.Categories)
	if err != nil {
		return MedicationWithModifications{}, err
	}

	return MedicationWithModifications{
		MedicationID:   src.MedicationID,
		MedicationName: src.MedicationName,
		Categories:     cats,
	}, nil
}

func MapToxicityWithGrades[T any](r T) (ToxicityWithGrades, error) {
	row := mapToToxicityWithGrades(r)
	var grades []ToxicityGrade
	if err := json.Unmarshal(row.Grades, &grades); err != nil {
		fmt.Printf("failed to decode grades for id %s: %v", row.ID, err)
		return ToxicityWithGrades{}, err
	}

	return_item := ToxicityWithGrades{
		ID:          row.ID,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
		Title:       row.Title,
		Category:    row.Category,
		Description: row.Description,
		Grades:      grades,
	}

	return return_item, nil
}

func MapToxicityWithGradesWithAdjust[T any](r T) (ToxicityWithGradesAndAdjustments, error) {
	row := mapToToxicityWithGrades(r)
	var grades []ToxicityGradeWithAdjustment
	if err := json.Unmarshal(row.Grades, &grades); err != nil {
		fmt.Printf("failed to decode grades for id %s: %v", row.ID, err)
		return ToxicityWithGradesAndAdjustments{}, err
	}

	return_item := ToxicityWithGradesAndAdjustments{
		ID:          row.ID,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
		Title:       row.Title,
		Category:    row.Category,
		Description: row.Description,
		Grades:      grades,
	}

	return return_item, nil
}

func MapToToxicityGrade(tox database.ToxicityGrade) ToxicityGrade {
	return ToxicityGrade{
		ID:          tox.ID,
		CreatedAt:   tox.CreatedAt,
		UpdatedAt:   tox.UpdatedAt,
		Description: tox.Description,
		Grade:       string(tox.Grade),
	}
}

type MedModificationResp struct {
	ID           string `json:"id"`
	MedicationID string `json:"medication_id"`

	Category    string `json:"category"`
	Subcategory string `json:"subcategory"`
	Adjustment  string `json:"adjustment"`
}

func MapMedModification(src database.MedicationModification) MedModificationResp {
	return MedModificationResp{
		ID:           src.ID.String(),
		MedicationID: src.MedicationID.String(),
		Category:     string(src.Category),
		Subcategory:  src.Subcategory,
		Adjustment:   src.Adjustment,
	}
}

func mapArticleRef(src database.ArticleReference) ArticleReference {

	return ArticleReference{
		ID:        src.ID,
		Title:     src.Title,
		Authors:   src.Authors,
		Journal:   src.Journal,
		Year:      src.Year,
		Pmid:      src.Pmid,
		Doi:       src.Doi,
		CreatedAt: src.CreatedAt,
		UpdatedAt: src.UpdatedAt,
	}
}

func mapPhysician(src database.Physician) Physician {
	return Physician{
		ID:        src.ID,
		FirstName: src.FirstName,
		LastName:  src.LastName,
		CreatedAt: src.CreatedAt,
		UpdatedAt: src.UpdatedAt,
	}
}

func mapSummaryProtocol(src database.Protocol) SummaryProtocol {
	return SummaryProtocol{
		ID:          src.ID,
		TumorGroup:  src.TumorGroup,
		Code:        src.Code,
		Name:        src.Name,
		Tags:        src.Tags,
		Notes:       src.Notes,
		ProtocolUrl: src.ProtocolUrl,
		HandOutUrl:  src.PatientHandoutUrl,
		CreatedAt:   src.CreatedAt,
		UpdatedAt:   src.UpdatedAt,
		RevisedOn:   src.RevisedOn,
		ActivatedOn: src.ActivatedOn,
	}
}

func MapEligibilityCriterion(src database.ProtocolEligibilityCriterium) ProtocolEligibilityCriterion {
	return ProtocolEligibilityCriterion{
		ID:          src.ID,
		Type:        src.Type,
		Description: src.Description,
		CreatedAt:   src.CreatedAt,
		UpdatedAt:   src.UpdatedAt,
	}
}

func MapPrecaution(src database.ProtocolPrecaution) ProtocolPrecaution {
	return ProtocolPrecaution{
		ID:          src.ID,
		Title:       src.Title,
		Description: src.Description,
		CreatedAt:   src.CreatedAt,
		UpdatedAt:   src.UpdatedAt,
	}
}

func MapCaution(src database.ProtocolCaution) ProtocolCaution {
	return ProtocolCaution{
		ID:          src.ID,
		Description: src.Description,
		CreatedAt:   src.CreatedAt,
		UpdatedAt:   src.UpdatedAt,
	}
}

func MapCycle(src database.ProtocolCycle) ProtocolCycle {
	return ProtocolCycle{
		ID:            src.ID,
		Cycle:         src.Cycle,
		CycleDuration: src.CycleDuration,
		CreatedAt:     src.CreatedAt,
		UpdatedAt:     src.UpdatedAt,
	}
}

func MapCycleWithTreatments(src database.ProtocolCycle) ProtocolCycle {
	return ProtocolCycle{
		ID:            src.ID,
		Cycle:         src.Cycle,
		CycleDuration: src.CycleDuration,
	}
}

func MapToMedicationModifications(rows []database.GetMedicationModificationsByProtocolRow) []MedicationModification {
	// Map to group medications by MedicationID
	medicationMap := make(map[uuid.UUID]*MedicationModification)

	for _, row := range rows {
		// Check if the medication already exists in the map
		if _, exists := medicationMap[row.MedicationID]; !exists {
			// Add a new medication to the map
			medicationMap[row.MedicationID] = &MedicationModification{
				MedicationID:         row.MedicationID,
				Medication:           row.Name,
				ModificationCategory: []ModificationCategory{},
			}
		}

		// Find or create the ModificationCategory within the Medication
		var category *ModificationCategory
		for i := range medicationMap[row.MedicationID].ModificationCategory {
			if medicationMap[row.MedicationID].ModificationCategory[i].Category == string(row.ModificationCategory) {
				category = &medicationMap[row.MedicationID].ModificationCategory[i]
				break
			}
		}

		// If the category doesn't exist, create it
		if category == nil {
			category = &ModificationCategory{
				Category:      string(row.ModificationCategory),
				Modifications: []Modifications{},
			}
			medicationMap[row.MedicationID].ModificationCategory = append(medicationMap[row.MedicationID].ModificationCategory, *category)
		}

		// Add the modification to the category
		category.Modifications = append(category.Modifications, Modifications{
			ID:          row.ModificationID,
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
