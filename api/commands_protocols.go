package api

import (
	"bcca_crawler/internal/config"
	"bcca_crawler/internal/database"	
	"github.com/google/uuid"
	"context"	
	"fmt"
	"strings"
	"encoding/json"	
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
	items,err := c.Db.GetArticleReferencesByProtocol(ctx,protocolID)
	if err != nil {
		return nil, err
	}

	returned_items := MapAll(items,mapArticleRef)	
	return returned_items, nil
}

func GetProtocolPhysicians(c *config.Config,ctx context.Context,protocolID uuid.UUID) ([]Physician, error) {
	items,err := c.Db.GetPhysicianByProtocol(ctx,protocolID)
	if err != nil {
		return nil, err
	}

	returned_items := MapAll(items,mapPhysician)	
	return returned_items, nil
}

func GetProtocolCautions(c *config.Config,ctx context.Context,protocolID uuid.UUID) ([]ProtocolCaution, error) {
	items,err := c.Db.GetProtocolCautionsByProtocol(ctx,protocolID)
	if err != nil {
		return nil, err
	}

	returned_items := MapAll(items,MapCaution)		
	return returned_items, nil
}

func GetProtocolPrecautions(c *config.Config,ctx context.Context,protocolID uuid.UUID) ([]ProtocolPrecaution, error) {
	items,err := c.Db.GetProtocolPrecautionsByProtocol(ctx,protocolID)
	if err != nil {
		return nil, err
	}

	returned_items := MapAll(items,MapPrecaution)	
	return returned_items, nil
}

func GetProtocolEligibilityCriteria(c *config.Config,ctx context.Context,protocolID uuid.UUID) ([]ProtocolEligibilityCriterion, error) {
	items,err := c.Db.GetEligibilityByProtocol(ctx,protocolID)
	if err != nil {
		return nil, err
	}

	returned_items := MapAll(items,MapEligibilityCriterion)	
	return returned_items, nil
}

func GetTests(c *config.Config,ctx context.Context,protocolID uuid.UUID) (LabsByProtocol, error) {
	items,err := c.Db.GetTestsByProtocol(ctx,protocolID)
	var labs LabsByProtocol
	if err != nil {
		return labs, err
	}

	if err := json.Unmarshal(items, &labs); err != nil {
		return labs, fmt.Errorf("failed to unmarshal tests json: %w", err)
	}

	return labs,nil  

}

func GetProtocolToxicities(c *config.Config,ctx context.Context,protocolID uuid.UUID) ([]ToxicityWithGradesAndAdjustments, error) {
	items,err := c.Db.GetToxicitiesWithGradesAndAdjustmentsByProtocol(ctx,protocolID)
	if err != nil {
		return nil, err
	}

	returned_value,err := MapAllWithError(items,MapToToxicityWithGradesAndAdjustmentsByProtocol)

	if err != nil {
		return nil, err
	}	
	
	return returned_value, nil
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


func CMD_GetProtocolBy(c *config.Config,ctx context.Context, ByWhat string, arg string) (ProtocolSumPayload, error) {
	
	Payload := ProtocolSumPayload{}
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

	r_tests,err := GetTests(c,ctx,protocol.ID)
	if err != nil {
		fmt.Println("Error getting tests: ", err)
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

	Payload = ProtocolSumPayload{
		ProtocolSummary:            r_protocol,
		ProtocolEligibilityCriteria: r_elig_criterias,
		ProtocolPrecautions:        r_precautions,
		ProtocolCautions:           r_cautions,
		Tests:                      r_tests,
		ProtocolCycles:             r_cycles,
		TreatmentModifications:     r_med_mod,
		Toxicities:     			r_tox_mod,
		Physicians:                 r_contact_physicians,
		ArticleReferences:          r_articles,
	}

	return Payload, nil
	
}