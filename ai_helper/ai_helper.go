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
	"github.com/lib/pq"
	"bcca_crawler/internal/config"
	"database/sql"

)

type Payload struct {
	ArticleReferences          []ArticleReference          `json:"ArticleReferences"`
	Medications                []Medication                `json:"Medications"`
	Physicians                 []Physician                 `json:"Physicians"`
	Protocols                  []Protocol                  `json:"Protocols"`
	ProtocolEligibilityCriteria []ProtocolEligibilityCriterion `json:"ProtocolEligibilityCriteria"`
	ProtocolPrecautions        []ProtocolPrecaution        `json:"ProtocolPrecautions"`
	ProtocolCautions		   []ProtocolCaution		   `json:"ProtocolCautions"`
	ProtocolPremedications     []ProtocolPremedication     `json:"ProtocolPremedications"`
	ProtocolSupportiveMedications []ProtocolSupportiveMedication `json:"ProtocolSupportiveMedications"`
	ProtocolCycles             []ProtocolCycle             `json:"ProtocolCycles"`
	Tests                      Tests                       `json:"Tests"`
	ToxicityModifications      []ToxicityModification      `json:"ToxicityModifications"`
}

type ArticleReference struct {
	Title   string `json:"Title"`
	Authors string `json:"Authors"`
	Journal string `json:"Journal"`
	Year    string `json:"Year"`
	Pmid    string `json:"Pmid"`
	Joi     string `json:"Joi"`
}

type Medication struct {
	Name        string `json:"Name"`
	Description string `json:"Description"`
	Category    string `json:"Category"`
}

type Physician struct {
	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
}

type Protocol struct {
	TumorGroup string `json:"TumorGroup"`
	Code       string `json:"Code"`
	Name       string `json:"Name"`
}

type ProtocolEligibilityCriterion struct {
	Type        string `json:"Type"`
	Description string `json:"Description"`
}

type ProtocolPrecaution struct {
	Title       string `json:"Title"`
	Description string `json:"Description"`
}

type ProtocolCaution struct {	
	Description string `json:"Description"`
}

type ProtocolPremedication struct {
	Medication string `json:"Medication"`
	Dose       string `json:"Dose"`
	Route      string `json:"Route"`
	Frequency  string `json:"Frequency"`
	Duration   string `json:"Duration"`
	Notes      string `json:"Notes"`
}

type ProtocolSupportiveMedication struct {
	Medication string `json:"Medication"`
	Dose       string `json:"Dose"`
	Route      string `json:"Route"`
	Frequency  string `json:"Frequency"`
	Duration   string `json:"Duration"`
	Notes      string `json:"Notes"`
}

type ProtocolCycle struct {
	Cycle         string         `json:"Cycle"`
	CycleDuration string         `json:"CycleDuration"`
	Treatments    []Treatment    `json:"Treatments"`
}

type Treatment struct {
	Medication            string                 `json:"Medication"`
	Dose                  string                 `json:"Dose"`
	Route                 string                 `json:"Route"`
	Frequency             string                 `json:"Frequency"`
	Duration              string                 `json:"Duration"`
	AdministrationGuide   string                 `json:"AdministrationGuide"`
	TreatmentModifications []TreatmentModification `json:"TreatmentModifications"`
}

type TreatmentModification struct {
	Category    string `json:"Category"`
	Description string `json:"Description"`
	Adjustement string `json:"Adjustement"`
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

type ToxicityModification struct {
	Title       string `json:"Title"`
	Grade       string `json:"Grade"`
	Adjustement string `json:"Adjustement"`
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
      "Joi": "string"
    }
  ],
  "Medications": [
    {      
      "Name": "string",
      "Description": "string (nullable)",
      "Category": "string"
    }
  ],
  "Physicians": [
    {      
      "FirstName": "string",
      "LastName": "string"
    }
  ],
  "Protocols": [
    {
      "TumorGroup": "string",
      "Code": "string",
      "Name": "string"
    }
  ],
  "ProtocolEligibilityCriteria": [
    {
      "Type": "string", // Inclusion, Exclusion or Other, Split each line as a separate object
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
  "ProtocolPremedications": [
    {
      "Medication": "string",
      "Dose": "string",
      "Route": "string",
      "Frequency": "string",
      "Duration": "string",
      "Notes": "string"
    }
  ],
  "ProtocolSupportiveMedications": [
    {
      "Medication": "string",
      "Dose": "string",
      "Route": "string",
      "Frequency": "string",
      "Duration": "string",
      "Notes": "string"
    }
  ],
  "ProtocolCycles": [
    {
	"Cycle": "string", // If there are no Cycles specified, mention "Cycle 1+"
	"CycleDuration": "string", //If not specified, mention "28 days"
	"Treatments": [
		{
		"Medication": "string",
		"Dose": "string",
		"Route": "string",
		"Frequency": "string",
		"Duration": "string",
		"AdministrationGuide": "string",
		"TreatmentModifications": [
			{			
			"Category": "string", // Hepatic or Renal Impairment
			"Description": "string", // Mild, Moderate, Severe
			"Adjustement": "string"			
			}
		]
		}
		]
	}
  ],
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
  "ToxicityModifications": [
    {      
      "Title": "string",
      "Grade": "string",
      "Adjustement": "string"
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

func (s *Session) AnalyzePDF(pdfURL string) (Payload, error) {
	pdfBytes, err := downloadPDF(pdfURL)
	if err != nil {
		return Payload{}, err
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

	fmt.Println("Payload: ", payload)
	protocol,err := s.Db.CreateProtocol(ctx, database.CreateProtocolParams{		
		TumorGroup: payload.Protocols[0].TumorGroup,
		Code: payload.Protocols[0].Code,
		Name: payload.Protocols[0].Name,
		Tags: []string{},
		Notes: payload.Protocols[0].Name,
	})
	if err != nil {		
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			// Duplicate key value violation
			fmt.Println("Record already exists")
			protocol,err = s.Db.GetProtocolByCode(ctx, payload.Protocols[0].Code)
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
			Joi: article.Joi,
		})
		if err != nil {		
			if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
				// Duplicate key value violation
				fmt.Println("Record already exists: ", article.Title)
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

		err = s.Db.AddArticleReferenceToProtocol(ctx, database.AddArticleReferenceToProtocolParams{
			ProtocolID: protocol.ID,
			ReferenceID: articleRef.ID,
		})
		if err != nil {
			fmt.Println("Error adding article reference to protocol: ", err)
		}			
	}

	

	// Create medications
	for _, medication := range payload.Medications {
		_,err := s.Db.AddMedication(ctx, database.AddMedicationParams{
			Name: medication.Name,
			Description: sql.NullString{String: medication.Description, Valid: medication.Description != ""},
			Category: medication.Category,
		})
		if err != nil {
			fmt.Println("Error in trying to create Med: ", medication.Name)
			fmt.Println("Error: ", err)
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
			Site: "Vancouver Centre",
		})
		if err != nil {
			if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
				// Duplicate key value violation
				fmt.Println("Physician Record already exists so retrieving from database.")
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

		err = s.Db.AddPhysicianToProtocol(ctx, database.AddPhysicianToProtocolParams{
			PhysicianID: phys.ID,
			ProtocolID: protocol.ID,
		})
		if err != nil {
			fmt.Printf("Error adding %v %v to protocol:%v\n",phys.FirstName,phys.LastName, err)
		}
			
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
				fmt.Println("Eligibility Record already exists so retrieving from database.")
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

		err = s.Db.LinkEligibilityToProtocol(ctx, database.LinkEligibilityToProtocolParams{
			ProtocolID: protocol.ID,
			CriteriaID: elig.ID,
		})
		if err != nil {
			fmt.Printf("Error adding %v:%v, to protocol:%v\n",elig.Type,elig.Description, err)
		}
			
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
				fmt.Println("Precaution Record already exists so retrieving from database.")
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

		err = s.Db.AddProtocolPrecautionToProtocol(ctx, database.AddProtocolPrecautionToProtocolParams{
			ProtocolID: protocol.ID,
			PrecautionID: precaut.ID,
		})

		if err != nil {
			fmt.Printf("Error adding %v:%v, to protocol:%v\n",precaut.Title,precaut.Description, err)
		}
	
	}


	// Create Protocol Cautions
	
	for _, caution := range payload.ProtocolCautions {
		caut,err := s.Db.CreateProtocolCaution(ctx, caution.Description)
		if err != nil {
			if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
				// Duplicate key value violation
				fmt.Println("Caution Record already exists so retrieving from database.")
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

		err = s.Db.AddProtocolCautionToProtocol(ctx, database.AddProtocolCautionToProtocolParams{
			ProtocolID: protocol.ID,
			CautionID: caut.ID,
		})

		if err != nil {
			fmt.Printf("Error adding caution: %v, to protocol:%v\n",caut.Description, err)
		}
	}
	

	for _, test := range payload.Tests.Baseline.RequiredBeforeTreatment {
		add,err := s.Db.AddTest(ctx, database.AddTestParams{
			Name: test,
			Description: sql.NullString{String: "", Valid: false},
		})
		if err != nil {
			if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
				// Duplicate key value violation
				fmt.Println("Test Record already exists so retrieving from database.")
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

		err = s.Db.AddBaselineTest(ctx, database.AddBaselineTestParams{
			ProtocolID: protocol.ID,
			TestID: add.ID,
		})

		if err != nil {
			fmt.Printf("Error adding baseline test: %v, to protocol:%v\n",add.Name, err)
		}

	}	

	// Create Tests, Non Urgent	

	for _, test := range payload.Tests.Baseline.RequiredButCanProceed {
		add,err := s.Db.AddTest(ctx, database.AddTestParams{
			Name: test,
			Description: sql.NullString{String: "", Valid: false},
		})
		if err != nil {
			if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
				// Duplicate key value violation
				fmt.Println("Test Record already exists so retrieving from database.")
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
		
		err = s.Db.AddNonUrgentTest(ctx, database.AddNonUrgentTestParams{
			ProtocolID: protocol.ID,
			TestID: add.ID,
		})

		if err != nil {
			fmt.Printf("Error adding baseline test non urgent: %v, to protocol:%v\n",add.Name, err)
		}

	}

	
	// Create Tests, If clinically indicated
	
	for _, test := range payload.Tests.Baseline.IfClinicallyIndicated {
		add,err := s.Db.AddTest(ctx, database.AddTestParams{
			Name: test,
			Description: sql.NullString{String: "", Valid: false},
		})
		if err != nil {
			if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
				// Duplicate key value violation
				fmt.Println("Test Record already exists so retrieving from database.")
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
		
		err = s.Db.AddIfNecessaryTest(ctx, database.AddIfNecessaryTestParams{
			ProtocolID: protocol.ID,
			TestID: add.ID,
		})

		if err != nil {
			fmt.Printf("Error adding baseline test if clinically indicated: %v, to protocol:%v\n",add.Name, err)
		}
	}	

	// Create Tests, follow-up
	
	for _, test := range payload.Tests.FollowUp.Required {
		add,err := s.Db.AddTest(ctx, database.AddTestParams{
			Name: test,
			Description: sql.NullString{String: "", Valid: false},
		})
		if err != nil {
			if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
				// Duplicate key value violation
				fmt.Println("Test Record already exists so retrieving from database.")
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
		
		err = s.Db.AddFollowupTest(ctx, database.AddFollowupTestParams{
			ProtocolID: protocol.ID,
			TestID: add.ID,
		})

		if err != nil {
			fmt.Printf("Error adding follow-up test required: %v, to protocol:%v\n",add.Name, err)
		}

	}

	// Create Tests, if necessary
	
	for _, test := range payload.Tests.FollowUp.IfClinicallyIndicated {
		add,err := s.Db.AddTest(ctx, database.AddTestParams{
			Name: test,
			Description: sql.NullString{String: "", Valid: false},
		})
		if err != nil {
			if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
				// Duplicate key value violation
				fmt.Println("Test Record already exists so retrieving from database.")
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
		
		err = s.Db.AddFollowupIfNecessaryTest(ctx, database.AddFollowupIfNecessaryTestParams{
			ProtocolID: protocol.ID,
			TestID: add.ID,
		})

		if err != nil {
			fmt.Printf("Error adding follow-up test if clinically indicated: %v, to protocol:%v\n",add.Name, err)
		}
		
	}
	

	// Create Toxicity Modifications
	
	for _, tox := range payload.ToxicityModifications {
		_,err := s.Db.AddToxicityModification(ctx, database.AddToxicityModificationParams{
			Title: tox.Title,
			Grade: tox.Grade,
			Adjustement: tox.Adjustement,
			ProtocolID: protocol.ID,
			
		})
		if err != nil {
			fmt.Println("Error creating toxicity modification: ", err)
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
			fmt.Println("Error creating protocol cycle: ", err)
		}

		for _, treatx := range cyc.Treatments {
			med,err := s.Db.AddMedication(ctx, database.AddMedicationParams{
				Name: treatx.Medication,
				Description: sql.NullString{String: "", Valid: false},
				Category: "Treatment",
			})
			if err != nil {

				if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
					// Duplicate key value violation
					fmt.Println("Medication Record already exists so retrieving from database.")
					med,err = s.Db.GetMedicationByName(ctx, treatx.Medication)

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

			for _, mod := range treatx.TreatmentModifications {
				_,err := s.Db.AddTreatmentModification(ctx, database.AddTreatmentModificationParams{
					Category: mod.Category,
					Description: mod.Description,
					Adjustement: mod.Adjustement,
					TreatmentID: treatment.ID,
				})
				if err != nil {
					fmt.Println("Error creating treatment modification: ", err)
				}
			}
		
			_err := s.Db.AddTreatmentToCycle(ctx, database.AddTreatmentToCycleParams{
				ProtocolCyclesID: cycle.ID,
				ProtocolTreatmentID: treatment.ID,				
			})
			if _err != nil {
				fmt.Println("Error adding treatment to cycle: ", _err)
			}
			

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
func handleRequest(ctx context.Context, model *genai.GenerativeModel, reqParts []genai.Part) (payload Payload, err error) {
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
			var payload Payload
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