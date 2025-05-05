package ai_helper

import (
	"context"
	"fmt"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
	"io"
	"net/http"
	"time"
	"strings"
	"encoding/json"
	"bcca_crawler/internal/database"
	"bcca_crawler/api"
	"github.com/lib/pq"
	"bcca_crawler/internal/config"	

)


type ProtocolPayload struct {	
	ProtocolSummary             api.SummaryProtocol                  `json:"SummaryProtocol"`
	Medications                []Medication               			 `json:"Medications"`
	ProtocolEligibilityCriteria []api.ProtocolEligibilityCriterion	 `json:"ProtocolEligibilityCriteria"`
	ProtocolPrecautions        []api.ProtocolPrecaution 		     `json:"ProtocolPrecautions"`
	ProtocolCautions		   []api.ProtocolCaution			 	 `json:"ProtocolCautions"`
	Tests                      api.Tests                     		 `json:"Tests"`
	ProtocolCycles             []api.ProtocolCycle           		 `json:"ProtocolCycles"`	
	Toxicities			       []api.Toxicity      				     `json:"Toxicities"`	
	Physicians                 []api.Physician                		 `json:"Physicians"`
	ArticleReferences          []api.ArticleReference        		 `json:"ArticleReferences"`
}

type Medication struct {
	Name        string `json:"Name"`
	Description string `json:"Description"`
	Category    string `json:"Category"`
	ModificationCategory 	[]api.ModificationCategory 			`json:"ModificationCategory"`
}

const ai_prompt = `You are an AI model tasked with analyzing a PDF document and extracting structured information in JSON format. The extracted data must strictly conform to the following JSON structure:

{
  "ArticleReferences": [
    {
      "Title": "string",
      "Authors": "string",
      "Journal": "string",
      "Year": "string",
      "Pmid": "string",
      "Doi": "string"
    }
  ],
  "Medications": [
    {      
      "Name": "string", //This should be a single word
      "Description": "string",
      "Category": "string" // Antibiotic, Antiemetic, etc.
	  "MedicationModifications": [
			{			
			"Category": "string", // Hepatic or Renal Impairment
			Modifications: [{			
				"Description": "string", // Mild, Moderate, Severe
				"Adjustement": "string"
				}]			
			}
    }
  ],
  "Physicians": [
    {      
      "FirstName": "string",
      "LastName": "string"
    }
  ],
  "SummaryProtocol":
    {
      "TumorGroup": "string",
      "Code": "string",
      "Name": "string"
	  "ActivatedOn": "string", // Date format: YYYY-MMM-DD
	  "RevisedOn": "string" // Date format: YYYY-MMM-DD
    }
  ,
  "ProtocolEligibilityCriteria": [
    {
      "Type": "string", // inclusion, exclusion or unknown, Split each line as a separate object
      "Description": "string"
    }
  ],
  "ProtocolPrecautions": [
    {
      "Title": "string",
      "Description": "string"
    }
  ],
  "ProtocolCautions": [
    {
	  "Description": "string"
	}
],
  "ProtocolCycles": [{
	"Cycle": "string", // If there are no Cycles specified, return "Cycle 1+"
	"CycleDuration": "string", //If blank, return "28 days"
	"Treatments": [{
		"Medication": "string",
		"Dose": "string",
		"Route": "string", // iv, oral, sc, im, topical, inhalation, other.
		"Frequency": "string", 
		"Duration": "string", 
		"AdministrationGuide": "string"
		}]
	}]
,
"Tests": {
  "Baseline": {
    "RequiredBeforeTreatment": ["string"],
    "RequiredButCanProceed": ["string"],
    "IfClinicallyIndicated": ["string"]
  },
  "FollowUp": {
    "Required": ["string"],
    "IfClinicallyIndicated": ["string"]
  }
},
  "Toxicities": [
    {      
      "Title": "string", // Neuropathy, Thrombopenia, Neutropenia, Diarrhea, etc.
      "Description": "string",
	  "Category": "string", // Hematologic, Neurologic, Gastrointestinal, etc.
	  "Modifications": [
	  				{
	  					"Grade": "string", // ONLY the grade number (1,2,3 or 4)
						"GradeDescription": "string", 
						"Adjustement": "string" // Dose reduction, Delay, Discontinuation
					}
				]     
    }
  ]
}

### Task:
1. Parse the provided PDF thoroughly.
2. Extract all relevant information and populate the JSON object as per the above structure.
3. Return the completed JSON object, ensuring all fields are validated for data type and consistency.`


type Session struct {
	ctx context.Context
	client  *genai.Client
	model   *genai.GenerativeModel
	
}

func NewSession(ctx context.Context, s *config.Config) (*Session, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(s.GeminiApiKey))
	if err != nil {
		return nil, fmt.Errorf("error creating AI client: %v", err)
	}

	model := client.GenerativeModel("gemini-1.5-flash")

	return &Session{
		ctx: ctx,
		client: client,
		model: model,
	}, nil
}

func (s *Session) Close() {	
	s.client.Close()
}

func (s *Session) AnalyzePDF(pdfURL string) (ProtocolPayload, error) {
	pdfBytes, err := downloadPDF(pdfURL)
	if err != nil {
		return ProtocolPayload{}, err
	}

	// Create the request.
	req := []genai.Part{
		genai.Blob{MIMEType: "application/pdf", Data: pdfBytes},
		genai.Text(ai_prompt),
	}

	return handleRequest(s.ctx, s.model, req)
}

func GetAiData(s *config.Config,proto string) error {
	ctx := context.Background()
	// Access your API key as an environment variable
	client, err := genai.NewClient(ctx, option.WithAPIKey(s.GeminiApiKey))
	if err != nil {
		return fmt.Errorf("error creating AI client: %v", err)
	}
	defer client.Close()

	done := make(chan bool)
	go spinner(done)

	model := client.GenerativeModel("gemini-1.5-flash")
	
	link := "http://www.bccancer.bc.ca/chemotherapy-protocols-site/Documents/Lymphoma-Myeloma/" + proto + "_Protocol.pdf"
 
	 pdfBytes, err := downloadPDF(link)
	 if err != nil {
		 return err
	 }
 
	 // Create the request.
	 req := []genai.Part{
		 genai.Blob{MIMEType: "application/pdf", Data: pdfBytes},
		 genai.Text(ai_prompt),		
	 }
 
	 payload, err := handleRequest(ctx, model, req)
	 if err != nil {
		return err
	 }
	fmt.Println("Payload:")
	api.PrintStruct(payload)
	protocol,err := s.Db.CreateProtocolbyScraping(ctx, database.CreateProtocolbyScrapingParams{		
		TumorGroup: payload.ProtocolSummary.TumorGroup,
		Code: payload.ProtocolSummary.Code,
		Name: payload.ProtocolSummary.Name,
		Tags: []string{},
		Notes: payload.ProtocolSummary.Name,
		RevisedOn: payload.ProtocolSummary.RevisedOn,
		ActivatedOn: payload.ProtocolSummary.ActivatedOn,
	})
	if err != nil {		
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			// Duplicate key value violation			
			protocol,err = s.Db.GetProtocolByCode(ctx, payload.ProtocolSummary.Code)
			if err != nil {
				fmt.Println("Error getting protocol: ", err)
				return err
			}
		} else {
			fmt.Println("Error creating protocol")
			return err
		}		
	}
	
	for _, article := range payload.ArticleReferences {
		articleRef,err := s.Db.CreateArticleReference(ctx, database.CreateArticleReferenceParams{
			Title: article.Title,
			Authors: article.Authors,
			Journal: article.Journal,
			Year: article.Year,
			Pmid: article.Pmid,
			Doi: article.Doi,
		})
		if err != nil {		
			if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
				// Duplicate key value violation				
				articleRef,err = s.Db.GetArticleReferenceByData(ctx, database.GetArticleReferenceByDataParams{
					Title: article.Title,
					Authors: article.Authors,
					Journal: article.Journal,
					Year: article.Year,
				})
				if err != nil {
					fmt.Println("Error getting protocol: ", err)
					return err
				}
			}		
		}

		_ = s.Db.AddArticleReferenceToProtocol(ctx, database.AddArticleReferenceToProtocolParams{
			ProtocolID: protocol.ID,
			ReferenceID: articleRef.ID,
		})				
	}

	

	// Create medications
	for _, medication := range payload.Medications {
		med,err := s.Db.AddMedication(ctx, database.AddMedicationParams{
			Name: medication.Name,
			Description: medication.Description,
			Category: medication.Category,
		})
		if err != nil {
			if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
				// Duplicate key value violation				
				med,err = s.Db.GetMedicationByName(ctx, medication.Name)
				if err != nil {
					fmt.Println("Error getting medication: ", err)
					return err
				}
			}			
		}
		
		for _, mod := range medication.ModificationCategory {
			for _, modx := range mod.Modifications {
				_,err := s.Db.AddMedicationModification(ctx, database.AddMedicationModificationParams{
					MedicationID: med.ID,
					Category: mod.Category,
					Subcategory: modx.Description,
					Adjustment: modx.Adjustment,
				})
				if err != nil {
					fmt.Println("Error creating medication modification: ", err)
				}
			}
		}
	}	

	// Create physicians	
	
	for _, physician := range payload.Physicians {
		email := strings.ToLower(physician.FirstName) + "." + strings.ToLower(physician.LastName) + "@bccancer.bc.ca"
		processedEmail := strings.ReplaceAll(email, " ", "")
		phys,err := s.Db.CreatePhysician(ctx, database.CreatePhysicianParams{
			FirstName: physician.FirstName,
			LastName: physician.LastName,
			Email: processedEmail,
			Site: "vancouver",
		})
		if err != nil {
			if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
				// Duplicate key value violation				
				phys,err = s.Db.GetPhysicianByName(ctx, database.GetPhysicianByNameParams{
					FirstName: physician.FirstName,
					LastName: physician.LastName,
				})
				if err != nil {
					fmt.Println("Error getting physician: ", err)
					return err
				}
			} else {
				fmt.Println("Error creating physician: ", err)
				return err
			}
		}

		_ = s.Db.AddPhysicianToProtocol(ctx, database.AddPhysicianToProtocolParams{
			PhysicianID: phys.ID,
			ProtocolID: protocol.ID,
		})		
			
	}	

	// Create protocol eligibility criteria
	
	for _, eligibility := range payload.ProtocolEligibilityCriteria {
		elig,err := s.Db.InsertEligibilityCriteria(ctx, database.InsertEligibilityCriteriaParams{
			Type: eligibility.Type,
			Description: eligibility.Description,			
		})
		if err != nil {
			if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
				// Duplicate key value violation
				elig,err = s.Db.GetElibilityCriteriaByDescription(ctx, eligibility.Description)
					
				if err != nil {
					fmt.Println("Error getting eligibility criteria: ", err)
					return err
				}
			} else {
				fmt.Println("Error creating eligibility criteria: ", err)
				return err
			}
		}

		_ = s.Db.LinkEligibilityToProtocol(ctx, database.LinkEligibilityToProtocolParams{
			ProtocolID: protocol.ID,
			CriteriaID: elig.ID,
		})	
			
	}	

	// Create protocol precautions
	
	for _, precaution := range payload.ProtocolPrecautions {
		precaut,err := s.Db.CreateProtocolPrecaution(ctx, database.CreateProtocolPrecautionParams{
			Title: precaution.Title,
			Description: precaution.Description,
		})
		if err != nil {
			if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
				// Duplicate key value violation
				precaut,err = s.Db.GetProtocolPrecautionByTitleAndDescription(ctx, database.GetProtocolPrecautionByTitleAndDescriptionParams{
					Title: precaution.Title,
					Description: precaution.Description,
				})
					
				if err != nil {
					fmt.Println("Error getting protocol precaution: ", err)
					return err
				}
			} else {
				fmt.Println("Error creating protocol precaution: ", err)
				return err
			}
		}

		_ = s.Db.AddProtocolPrecautionToProtocol(ctx, database.AddProtocolPrecautionToProtocolParams{
			ProtocolID: protocol.ID,
			PrecautionID: precaut.ID,
		})			
	}

	// Create Protocol Cautions
	
	for _, caution := range payload.ProtocolCautions {
		caut,err := s.Db.CreateProtocolCaution(ctx, caution.Description)
		if err != nil {
			if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
				// Duplicate key value violation				
				caut,err = s.Db.GetProtocolCautionByDescription(ctx, caution.Description)
					
				if err != nil {
					fmt.Println("Error getting protocol caution: ", err)
					return err
				}
			} else {
				fmt.Println("Error creating protocol caution: ", err)
				return err
			}
		}

		_ = s.Db.AddProtocolCautionToProtocol(ctx, database.AddProtocolCautionToProtocolParams{
			ProtocolID: protocol.ID,
			CautionID: caut.ID,
		})		
	}
	

	for _, test := range payload.Tests.Baseline.RequiredBeforeTreatment {
		add,err := s.Db.AddTest(ctx, database.AddTestParams{
			Name: test,
		})
		if err != nil {
			if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
				// Duplicate key value violation
				add,err = s.Db.GetTestByName(ctx, test)
					
				if err != nil {
					fmt.Println("Error getting test: ", err)
					return err
				}
			} else {
				fmt.Println("Error creating test: ", err)
				return err
			}
		}

		_,_ = s.Db.AddTestToProtocolByCategoryAndUrgency(ctx, database.AddTestToProtocolByCategoryAndUrgencyParams{
			ProtocolID: protocol.ID,
			TestID: add.ID,
			Category: database.CategoryEnumBaseline,
			Urgency: database.UrgencyEnumUrgent,
		})


	}	

	// Create Tests, Non Urgent	

	for _, test := range payload.Tests.Baseline.RequiredButCanProceed {
		add,err := s.Db.AddTest(ctx, database.AddTestParams{
			Name: test,
		})
		if err != nil {
			if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
				// Duplicate key value violation
				add,err = s.Db.GetTestByName(ctx, test)
					
				if err != nil {
					fmt.Println("Error getting test: ", err)
					return err
				}
			} else {
				fmt.Println("Error creating test: ", err)
				return err
			}
		}
		
		_,_ = s.Db.AddTestToProtocolByCategoryAndUrgency(ctx, database.AddTestToProtocolByCategoryAndUrgencyParams{
			ProtocolID: protocol.ID,
			TestID: add.ID,
			Category: database.CategoryEnumBaseline,
			Urgency: database.UrgencyEnumNonUrgent,
		})

	}

	
	// Create Tests, If clinically indicated
	
	for _, test := range payload.Tests.Baseline.IfClinicallyIndicated {
		add,err := s.Db.AddTest(ctx, database.AddTestParams{
			Name: test,
		})
		if err != nil {
			if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
				// Duplicate key value violation
				add,err = s.Db.GetTestByName(ctx, test)
					
				if err != nil {
					fmt.Println("Error getting test: ", err)
					return err
				}
			} else {
				fmt.Println("Error creating test: ", err)
				return err
			}
		}
		
		_,_ = s.Db.AddTestToProtocolByCategoryAndUrgency(ctx, database.AddTestToProtocolByCategoryAndUrgencyParams{
			ProtocolID: protocol.ID,
			TestID: add.ID,
			Category: database.CategoryEnumBaseline,
			Urgency: database.UrgencyEnumIfNecessary,
		})

	}	

	// Create Tests, follow-up
	
	for _, test := range payload.Tests.FollowUp.Required {
		add,err := s.Db.AddTest(ctx, database.AddTestParams{
			Name: test,
		})
		if err != nil {
			if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
				// Duplicate key value violation
				add,err = s.Db.GetTestByName(ctx, test)
					
				if err != nil {
					fmt.Println("Error getting test: ", err)
					return err
				}
			} else {
				fmt.Println("Error creating test: ", err)
				return err
			}
		}
		
		_,_ = s.Db.AddTestToProtocolByCategoryAndUrgency(ctx, database.AddTestToProtocolByCategoryAndUrgencyParams{
			ProtocolID: protocol.ID,
			TestID: add.ID,
			Category: database.CategoryEnumFollowup,
			Urgency: database.UrgencyEnumUrgent,
		})

	}

	// Create Tests, if necessary
	
	for _, test := range payload.Tests.FollowUp.IfClinicallyIndicated {
		add,err := s.Db.AddTest(ctx, database.AddTestParams{
			Name: test,
		})
		if err != nil {
			if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
				// Duplicate key value violation
				add,err = s.Db.GetTestByName(ctx, test)
					
				if err != nil {
					fmt.Println("Error getting test: ", err)
					return err
				}
			} else {
				fmt.Println("Error creating test: ", err)
				return err
			}
		}
		
		_,_ = s.Db.AddTestToProtocolByCategoryAndUrgency(ctx, database.AddTestToProtocolByCategoryAndUrgencyParams{
			ProtocolID: protocol.ID,
			TestID: add.ID,
			Category: database.CategoryEnumFollowup,
			Urgency: database.UrgencyEnumIfNecessary,
		})
		
	}
	

	// Create Toxicity Modifications
	
	for _, tox := range payload.Toxicities {

		toxicity,err := s.Db.AddToxicity(ctx, database.AddToxicityParams{
			Title: tox.Title,
			Description: tox.Description,
			Category: tox.Category,
		})

		if err != nil {
			if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
				// Duplicate key value violation
				toxicity,err = s.Db.GetToxicityByName(ctx, tox.Title)
					
				if err != nil {
					fmt.Println("Error getting toxicity: ", err)
					return err
				}
			} else {
				fmt.Println("Error creating toxicity: ", err)
				return err
			}
		}

		for _, mod := range tox.Modifications {
			grade,err := s.Db.AddToxicityGrade(ctx, database.AddToxicityGradeParams{
				Grade: database.GradeEnum(mod.Grade),
				Description: mod.GradeDescription,
				ToxicityID: toxicity.ID,
			})
			if err != nil {
				if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
					// Duplicate key value violation
					grade,err = s.Db.GetToxicityGradeByGrade(ctx, database.GetToxicityGradeByGradeParams{
						Grade: database.GradeEnum(mod.Grade),
						ToxicityID: toxicity.ID,
					})
						
					if err != nil {
						fmt.Println("Error getting grade: ", err)
						return err
					}
				} else {
					fmt.Println("Error creating grade: ", err)
					return err
				}
			}

			_,err = s.Db.AddToxicityModification(ctx, database.AddToxicityModificationParams{
				Adjustment: mod.Adjustment,
				ToxicityGradeID: grade.ID,
				ProtocolID: protocol.ID,
			})

			if err != nil {
				fmt.Println("Error creating toxicity modification: ", err)
			}
		}

	}

	
	// Create Protocol Premedications
	// Create Protocol Supportive Medications

	// Create Protocol Cycles
	for _, cyc := range payload.ProtocolCycles {
		cycle,err := s.Db.AddCycleToProtocol(ctx, database.AddCycleToProtocolParams{
			Cycle: cyc.Cycle,
			CycleDuration: cyc.CycleDuration,
			ProtocolID: protocol.ID,
		})
		if err != nil {
			if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
				// Duplicate key value violation
				fmt.Println("Cycle Record already exists so retrieving from database.")
				cycle,err = s.Db.GetCycleByData(ctx, database.GetCycleByDataParams{
					Cycle: cyc.Cycle,
					CycleDuration: cyc.CycleDuration,
					ProtocolID: protocol.ID,
				})

				if err != nil {
					fmt.Println("Error getting cycle: ", err)
					return err
				}
			} else {
				fmt.Println("Error creating cycle: ", err)
				return err
			}
			
		}

		for _, treatx := range cyc.Treatments {
			med,err := s.Db.AddMedication(ctx, database.AddMedicationParams{
				Name: treatx.MedicationName,
				Description: "",
				Category: "Treatment",
			})
			if err != nil {

				if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
					// Duplicate key value violation
					fmt.Println("Medication Record already exists so retrieving from database.")
					med,err = s.Db.GetMedicationByName(ctx, treatx.MedicationName)

					if err != nil {
						fmt.Println("Error getting medication: ", err)
						return err
					}
				} else {
					fmt.Println("Error creating medication: ", err)
					return err
				}
			}		

			treatment,err := s.Db.AddProtocolTreatment(ctx, database.AddProtocolTreatmentParams{
				Medication: med.ID,
				Dose: treatx.Dose,
				Route: treatx.Route,
				Frequency: treatx.Frequency,
				Duration: treatx.Duration,
				AdministrationGuide: treatx.AdministrationGuide,				
			})

			if err != nil {
				if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
					// Duplicate key value violation
					fmt.Println("Treatment Record already exists so retrieving from database.")
					treatment,err = s.Db.GetProtocolTreatmentByData(ctx, database.GetProtocolTreatmentByDataParams{
						Medication: med.ID,
						Dose: treatx.Dose,
						Route: treatx.Route,
						Frequency: treatx.Frequency,
						Duration: treatx.Duration,
					})

					if err != nil {
						fmt.Println("Error getting treatment: ", err)
						return err
					}
				} else {
					fmt.Println("Error creating treatment: ", err)
					return err
				}
			}			
		
			_ = s.Db.AddTreatmentToCycle(ctx, database.AddTreatmentToCycleParams{
				ProtocolCyclesID: cycle.ID,
				ProtocolTreatmentID: treatment.ID,				
			})
		}
	
	}	

	done <- true  
	return nil
}


func spinner(done chan bool) {
	// Spinner characters
	chars := []rune{'|', '/', '-', '\\'}
	i := 0
	timelapsed := 0.0

	for {
		select {
		case <-done:
			// Exit the spinner loop when the done channel is triggered
			fmt.Printf("Done! Time elapsed: %.2fs\n", timelapsed)			
			return
		default:
			// Display the spinner
			fmt.Printf("\r%c", chars[i])
			i = (i + 1) % len(chars)
			time.Sleep(100 * time.Millisecond) // Adjust for spinner speed
			timelapsed += 0.1
		}
	}
}

// Function to download PDF content
func downloadPDF(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// Function to extract JSON data from the string content between backticks
func extractJSON(contentStr string) (string, error) {
	start := strings.Index(contentStr, "```json")
	if start == -1 {
		return "", fmt.Errorf("no opening backticks found")
	}

	end := strings.Index(contentStr[start+7:], "```")
	if end == -1 {
		return "", fmt.Errorf("no closing backticks found")
	}

	return contentStr[start+7 : start+7+end], nil
}

// Function to handle each request, send it, and parse the response
func handleRequest(ctx context.Context, model *genai.GenerativeModel, reqParts []genai.Part) (payload ProtocolPayload, err error) {
	resp, err := model.GenerateContent(ctx, reqParts...)
	if err != nil {
		fmt.Printf("Request failed: %v\n", err)
		return
	}

	for _, c := range resp.Candidates {
		if c.Content != nil {
			contentStr := fmt.Sprintf("%s", *c.Content)			

			// Extract content between backticks and unmarshal
			jsonData, err := extractJSON(contentStr)
			if err != nil {
				fmt.Printf("Request - Error extracting JSON: %v\n", err)
				continue
			}

			// Parse the extracted JSON data
			var payload ProtocolPayload
			err = json.Unmarshal([]byte(jsonData), &payload)
			if err != nil {
				fmt.Printf("Request - Error unmarshaling JSON: %v\n", err)
				continue
			}

			return payload, nil
		}
	}

	return payload, fmt.Errorf("no content generated for request %v", err)
}