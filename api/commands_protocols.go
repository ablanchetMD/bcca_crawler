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


func CMD_GetProtocolBy(c *config.Config,ByWhat string, arg string) (ProtocolPayload, error) {
	ctx := context.Background()

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

	articles,err := c.Db.GetArticleReferencesByProtocol(ctx,protocol.ID)
	if err != nil {
		fmt.Println("Error getting articles: ", err)
		return Payload, err
	}

	r_articles := []ArticleReference{}
	for _, a := range articles {
		r_articles = append(r_articles, mapArticleRef(a))
	}

	elig_criterias,err := c.Db.GetEligibilityByProtocol(ctx,protocol.ID)
	if err != nil {
		fmt.Println("Error getting eligibility criteria: ", err)
		return Payload, err
	}

	r_elig_criterias := []ProtocolEligibilityCriterion{}
	for _, ec := range elig_criterias {		
		r_elig_criterias = append(r_elig_criterias, mapEligibilityCriterion(ec))
	}

	
	baseline_tests,err := c.Db.GetBaselineTestsByProtocol(ctx,protocol.ID)

	if err != nil {
		fmt.Println("Error getting baseline tests: ", err)
		return Payload, err
	}

	r_baseline_tests := mapTest(baseline_tests)

	nonurgent_tests,err := c.Db.GetNonUrgentTestsByProtocol(ctx,protocol.ID)

	if err != nil {
		fmt.Println("Error getting nonurgent tests: ", err)
		return Payload, err
	}

	r_nonurgent_tests := mapTest(nonurgent_tests)

	ifnec_tests,err := c.Db.GetIfNecessaryTestsByProtocol(ctx,protocol.ID)

	if err != nil {
		fmt.Println("Error getting if necessary tests: ", err)
		return Payload, err
	}

	r_ifnec_tests := mapTest(ifnec_tests)

	followup_tests,err := c.Db.GetFollowupTestsByProtocol(ctx,protocol.ID)

	if err != nil {
		fmt.Println("Error getting followup tests: ", err)
		return Payload, err
	}

	r_followup_tests := mapTest(followup_tests)

	followup_ifnec_tests,err := c.Db.GetFollowupIfNecessaryTestsByProtocol(ctx,protocol.ID)

	if err != nil {
		fmt.Println("Error getting followup if necessary tests: ", err)
		return Payload, err
	}

	r_followup_ifnec_tests := mapTest(followup_ifnec_tests)

	contact_physicians,err := c.Db.GetPhysicianByProtocol(ctx,protocol.ID)


	if err != nil {
		fmt.Println("Error getting contact physicians: ", err)
		return Payload, err
	}

	r_contact_physicians := []Physician{}
	for _, p := range contact_physicians {
		r_contact_physicians = append(r_contact_physicians, mapPhysician(p))
	}

	cautions,err := c.Db.GetProtocolCautionsByProtocol(ctx,protocol.ID)

	if err != nil {
		fmt.Println("Error getting protocol cautions: ", err)
		return Payload, err
	}

	r_cautions := []ProtocolCaution{}
	for _, c := range cautions {
		r_cautions = append(r_cautions, mapCaution(c))
	}
	
	precautions,err := c.Db.GetProtocolPrecautionsByProtocol(ctx,protocol.ID)

	if err != nil {
		fmt.Println("Error getting protocol precautions: ", err)
		return Payload, err
	}

	r_precautions := []ProtocolPrecaution{}
	for _, p := range precautions {
		r_precautions = append(r_precautions, mapPrecaution(p))
	}

	cycles,err := c.Db.GetCyclesByProtocol(ctx,protocol.ID)

	if err != nil {
		fmt.Println("Error getting protocol cycles: ", err)
		return Payload, err
	}

	r_cycles := []ProtocolCycle{}
	

	for _, cyc := range cycles {
		r_cycle := mapCycle(cyc)
		treatments,err := c.Db.GetTreatmentsByCycle(ctx,cyc.ID)
		if err != nil {
			fmt.Println("Error getting treatments: ", err)
			return Payload, err
		}

		r_treatments := []Treatment{}
		for _, treat := range treatments {
			med,err := c.Db.GetMedicationByID(ctx,treat.Medication)
			if err != nil {
				fmt.Println("Error getting medication: ", err)
				return Payload, err
			}
			
			r_treatment := mapTreatment(treat)
			r_treatment.MedicationName = med.Name
			adjust,err := c.Db.GetTreatmentModificationsByTreatment(ctx,treat.ID)
			if err != nil {
				return Payload, err
			}

			r_adjust := []TreatmentModification{}
			for _, adj := range adjust {
				r_adjust = append(r_adjust, mapTreatmentModification(adj))
			}

			r_treatment.TreatmentModifications = r_adjust
			r_treatments = append(r_treatments, r_treatment)

		}

		r_cycle.Treatments = r_treatments
		r_cycles = append(r_cycles, r_cycle)
	}

	tox_mod, err := c.Db.GetToxicityModificationsByProtocol(ctx,protocol.ID)

	r_tox_mod := []ToxicityModification{}

	for _, tox := range tox_mod {
		r_tox_mod = append(r_tox_mod, mapToxicityModification(tox))
	}

	if err != nil {
		fmt.Println("Error getting toxicity modifications: ", err)
		return Payload, err
	}
	Payload = ProtocolPayload{
		ProtocolSummary:          r_protocol,
		ProtocolEligibilityCriteria: r_elig_criterias,
		ProtocolPrecautions:        r_precautions,
		ProtocolCautions:           r_cautions,
		Tests:                      Tests{
			Baseline: BaselineTests{
				RequiredBeforeTreatment: r_baseline_tests,
				RequiredButCanProceed:   r_nonurgent_tests,
				IfClinicallyIndicated:   r_ifnec_tests,
			},
			FollowUp: FollowUpTests{
				Required:               r_followup_tests,
				IfClinicallyIndicated:  r_followup_ifnec_tests,			
			
		},
		},
		ProtocolCycles:             r_cycles,
		ToxicityModifications:      r_tox_mod,
		Physicians:                 r_contact_physicians,
		ArticleReferences:          r_articles,
	}

	return Payload, nil
	
}