package api

import (
	"bcca_crawler/internal/config"
	"bcca_crawler/internal/database"
	"github.com/google/uuid"
	"context"	
	"fmt"
	"strings"
)


func CMD_DeleteProtocol(c *config.Config, arg string) error {
	ctx := context.Background()

	protocol,err := c.Db.GetProtocolByCode(ctx,arg)
	if err != nil {
		return err
	}
	
	err = c.Db.DeleteProtocol(ctx,protocol.ID)
	if err != nil {
		return err
	}
	fmt.Printf("Protocol %s deleted", arg)
	return nil
}

func GetProtocolReferences(c *config.Config,ctx context.Context,protocolID uuid.UUID) ([]ArticleReference, error) {
	articles,err := c.Db.GetArticleReferencesByProtocol(ctx,protocolID)
	if err != nil {
		return nil, err
	}

	r_articles := []ArticleReference{}
	for _, a := range articles {
		r_articles = append(r_articles, mapArticleRef(a))
	}
	return r_articles, nil
}

func GetProtocolPhysicians(c *config.Config,ctx context.Context,protocolID uuid.UUID) ([]Physician, error) {
	contact_physicians,err := c.Db.GetPhysicianByProtocol(ctx,protocolID)
	if err != nil {
		return nil, err
	}

	r_contact_physicians := []Physician{}
	for _, p := range contact_physicians {
		r_contact_physicians = append(r_contact_physicians, mapPhysician(p))
	}
	return r_contact_physicians, nil
}

func GetProtocolCautions(c *config.Config,ctx context.Context,protocolID uuid.UUID) ([]ProtocolCaution, error) {
	cautions,err := c.Db.GetProtocolCautionsByProtocol(ctx,protocolID)
	if err != nil {
		return nil, err
	}

	r_cautions := []ProtocolCaution{}
	for _, c := range cautions {
		r_cautions = append(r_cautions, MapCaution(c))
	}
	return r_cautions, nil
}

func GetProtocolPrecautions(c *config.Config,ctx context.Context,protocolID uuid.UUID) ([]ProtocolPrecaution, error) {
	precautions,err := c.Db.GetProtocolPrecautionsByProtocol(ctx,protocolID)
	if err != nil {
		return nil, err
	}

	r_precautions := []ProtocolPrecaution{}
	for _, p := range precautions {
		r_precautions = append(r_precautions, MapPrecaution(p))
	}
	return r_precautions, nil
}

func GetProtocolEligibilityCriteria(c *config.Config,ctx context.Context,protocolID uuid.UUID) ([]ProtocolEligibilityCriterion, error) {
	elig_criterias,err := c.Db.GetEligibilityByProtocol(ctx,protocolID)
	if err != nil {
		return nil, err
	}

	r_elig_criterias := []ProtocolEligibilityCriterion{}
	for _, ec := range elig_criterias {		
		r_elig_criterias = append(r_elig_criterias, MapEligibilityCriterion(ec))
	}
	return r_elig_criterias, nil
}

func GetBaselineTests(c *config.Config,ctx context.Context,protocolID uuid.UUID) (BaselineTests, error) {
	baseline_tests,err := c.Db.GetTestsByProtocolByCategoryAndUrgency(ctx,database.GetTestsByProtocolByCategoryAndUrgencyParams{
		ProtocolID: protocolID,
		Category: database.CategoryEnumBaseline,
		Urgency: database.UrgencyEnumUrgent,
	})
	if err != nil {
		return BaselineTests{}, err
	}

	r_baseline_tests := mapTest(baseline_tests)

	nonurgent_tests,err := c.Db.GetTestsByProtocolByCategoryAndUrgency(ctx,database.GetTestsByProtocolByCategoryAndUrgencyParams{
		ProtocolID: protocolID,
		Category: database.CategoryEnumBaseline,
		Urgency: database.UrgencyEnumNonUrgent,
	})
	if err != nil {
		return BaselineTests{}, err
	}

	r_nonurgent_tests := mapTest(nonurgent_tests)

	ifnec_tests,err := c.Db.GetTestsByProtocolByCategoryAndUrgency(ctx,database.GetTestsByProtocolByCategoryAndUrgencyParams{
		ProtocolID: protocolID,
		Category: database.CategoryEnumBaseline,
		Urgency: database.UrgencyEnumIfNecessary,
	})
	if err != nil {
		return BaselineTests{}, err
	}

	r_ifnec_tests := mapTest(ifnec_tests)

	return BaselineTests{
		RequiredBeforeTreatment: r_baseline_tests,
		RequiredButCanProceed:   r_nonurgent_tests,
		IfClinicallyIndicated:   r_ifnec_tests,
	}, nil
}

func GetFollowUpTests(c *config.Config,ctx context.Context,protocolID uuid.UUID) (FollowUpTests, error) {
	followup_tests,err := c.Db.GetTestsByProtocolByCategoryAndUrgency(ctx,database.GetTestsByProtocolByCategoryAndUrgencyParams{
		ProtocolID: protocolID,
		Category: database.CategoryEnumFollowup,
		Urgency: database.UrgencyEnumUrgent,
	})
	if err != nil {
		return FollowUpTests{}, err
	}

	r_followup_tests := mapTest(followup_tests)

	followup_ifnec_tests,err := c.Db.GetTestsByProtocolByCategoryAndUrgency(ctx,database.GetTestsByProtocolByCategoryAndUrgencyParams{
		ProtocolID: protocolID,
		Category: database.CategoryEnumFollowup,
		Urgency: database.UrgencyEnumIfNecessary,
	})
	if err != nil {
		return FollowUpTests{}, err
	}

	r_followup_ifnec_tests := mapTest(followup_ifnec_tests)

	return FollowUpTests{
		Required:               r_followup_tests,
		IfClinicallyIndicated:  r_followup_ifnec_tests,
	}, nil
}


func GetProtocolToxicities(c *config.Config,ctx context.Context,protocolID uuid.UUID) ([]Toxicity, error) {
	tox_mod,err := c.Db.GetToxicityModificationByProtocol(ctx,protocolID)
	if err != nil {
		return nil, err
	}
	
	return mapToToxicities(tox_mod), nil
}

func GetProtocolModifications(c *config.Config,ctx context.Context,protocolID uuid.UUID) ([]MedicationModification, error) {
	med_mod,err := c.Db.GetMedicationModificationsByProtocol(ctx,protocolID)
	if err != nil {
		return nil, err
	}
	
	return MapToMedicationModifications(med_mod), nil
}

func GetProtocolCycles(c *config.Config,ctx context.Context,protocolID uuid.UUID) ([]ProtocolCycle, error) {
	cycles,err := c.Db.GetCyclesByProtocol(ctx,protocolID)
	if err != nil {
		return nil, err
	}

	r_cycles := []ProtocolCycle{}
	for _, cyc := range cycles {
		r_cycle := MapCycle(cyc)
		treatments,err := c.Db.GetTreatmentsByCycle(ctx,cyc.ID)
		if err != nil {
			return nil, err
		}

		r_treatments := []Treatment{}
		for _, treat := range treatments {			
			
			r_treatment := MapTreatmentByCycle(treat)			
			r_treatments = append(r_treatments, r_treatment)

		}

		r_cycle.Treatments = r_treatments
		r_cycles = append(r_cycles, r_cycle)
	}

	return r_cycles, nil
}


func CMD_GetProtocolBy(c *config.Config,ctx context.Context, ByWhat string, arg string) (ProtocolPayload, error) {
	
	Payload := ProtocolPayload{}
	protocol := database.Protocol{}
	err := error(nil)

	var id uuid.UUID

	if strings.ToLower(ByWhat) == "id" {
		id,err = uuid.Parse(arg)
		if err != nil {
			return Payload, err
		}
	}

	switch strings.ToLower(ByWhat) {
		case "code":
			protocol,err = c.Db.GetProtocolByCode(ctx,arg)
		
		case "id":
			
			protocol,err = c.Db.GetProtocolByID(ctx,id)
	}	
	
	if err != nil {
		fmt.Println("Error getting protocol: ", err)
		return Payload, err
	}
	r_protocol := mapSummaryProtocol(protocol)

	r_articles,err := GetProtocolReferences(c,ctx,protocol.ID)
	if err != nil {
		fmt.Println("Error getting article references: ", err)
		return Payload, err
	}

	r_elig_criterias,err := GetProtocolEligibilityCriteria(c,ctx,protocol.ID)
	if err != nil {
		fmt.Println("Error getting eligibility criteria: ", err)
		return Payload, err
	}

	r_baseline_tests,err := GetBaselineTests(c,ctx,protocol.ID)
	if err != nil {
		fmt.Println("Error getting baseline tests: ", err)
		return Payload, err
	}

	r_followup_tests,err := GetFollowUpTests(c,ctx,protocol.ID)
	if err != nil {
		fmt.Println("Error getting followup tests: ", err)
		return Payload, err
	}

	r_contact_physicians,err := GetProtocolPhysicians(c,ctx,protocol.ID)
	if err != nil {
		fmt.Println("Error getting contact physicians: ", err)
		return Payload, err
	}

	r_cautions,err := GetProtocolCautions(c,ctx,protocol.ID)
	if err != nil {
		fmt.Println("Error getting protocol cautions: ", err)
		return Payload, err
	}

	r_precautions,err := GetProtocolPrecautions(c,ctx,protocol.ID)
	if err != nil {
		fmt.Println("Error getting protocol precautions: ", err)
		return Payload, err
	}	

	r_cycles,err := GetProtocolCycles(c,ctx,protocol.ID)

	if err != nil {
		fmt.Println("Error getting protocol cycles: ", err)
		return Payload, err
	}

	r_med_mod,err := GetProtocolModifications(c,ctx,protocol.ID)
	if err != nil {
		fmt.Println("Error getting medication modifications: ", err)
		return Payload, err
	}

	
	r_tox_mod,err := GetProtocolToxicities(c,ctx,protocol.ID)

	if err != nil {
		fmt.Println("Error getting toxicity modifications: ", err)
		return Payload, err
	}

	Payload = ProtocolPayload{
		ProtocolSummary:            r_protocol,
		ProtocolEligibilityCriteria: r_elig_criterias,
		ProtocolPrecautions:        r_precautions,
		ProtocolCautions:           r_cautions,
		Tests:                      Tests{
			Baseline: r_baseline_tests,
			FollowUp: r_followup_tests,
		},
		ProtocolCycles:             r_cycles,
		TreatmentModifications:     r_med_mod,
		Toxicities:     			r_tox_mod,
		Physicians:                 r_contact_physicians,
		ArticleReferences:          r_articles,
	}

	return Payload, nil
	
}