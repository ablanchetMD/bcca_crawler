package ai_helper

import (
	"bcca_crawler/api"
	"bcca_crawler/internal/config"

	"bcca_crawler/internal/database"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/lib/pq"
	"google.golang.org/genai"
)

type ProtocolPayload struct {
	ProtocolSummary             api.SummaryProtocol                `json:"summary_protocol"`
	ProtocolEligibilityCriteria []api.ProtocolEligibilityCriterion `json:"protocol_eligibility_criteria"`
	ProtocolPrecautions         []api.ProtocolPrecaution           `json:"protocol_precautions"`
	ProtocolCautions            []api.ProtocolCaution              `json:"protocol_cautions"`
	TestGroups                  []api.TestGroup                    `json:"test_groups"`
	PrescriptionGroups          []api.PrescriptionGroup            `json:"prescription_groups"`
	ProtocolCycles              []api.ProtocolCycle                `json:"protocol_cycles"`
	Toxicities                  []api.Toxicity                     `json:"toxicities"`
	Physicians                  []api.Physician                    `json:"physicians"`
	ArticleReferences           []api.ArticleReference             `json:"article_references"`
}

type Medication struct {
	Name                 string                     `json:"name"`
	Description          string                     `json:"description"`
	Category             string                     `json:"category"`
	AlternateNames       []string                   `json:"alternate_names"` // If you know the alternate names for this drug, include them here
	ModificationCategory []api.ModificationCategory `json:"modification_category"`
}

const ai_prompt = `You are a medical oncologist tasked with analyzing the joined PDF document and extracting structured information in JSON format.
### Task:
1. Parse the provided PDF thoroughly.
2. Extract all relevant information, and complete the JSON object according to the provided schema.
3. Ensure that the extracted information is accurate and complete it to the best of your expertise knowledge.
3. Toxicities should be defined using the CTCAE v5 terminology. Generate only a toxicity with adjustment if there are suggested guidances.
4. Each tests should be a single entity.
5. Return the completed JSON object, ensuring all fields are validated for data type and consistency.`

type Session struct {
	ctx    context.Context
	client *genai.Client
	model  string
}

func NewSession(ctx context.Context, s *config.Config) (*Session, error) {
	client, err := genai.NewClient(ctx,
		&genai.ClientConfig{
			APIKey:  s.GeminiApiKey,
			Backend: genai.BackendGeminiAPI,
		})
	if err != nil {
		return nil, fmt.Errorf("error creating AI client: %v", err)
	}

	model := "gemini-2.5-flash"

	return &Session{
		ctx:    ctx,
		client: client,
		model:  model,
	}, nil
}

func retry[T any](attempts int, sleep time.Duration, fn func() (T, error)) (T, error) {
	var zero T
	for i := 0; i < attempts; i++ {
		result, err := fn()
		if err == nil {
			return result, nil
		}
		time.Sleep(sleep * time.Duration(i+1)) // linear backoff
	}
	return zero, fmt.Errorf("after %d attempts, failed", attempts)
}

func RunAllLinks(s *config.Config, links []string) {
	ctx := context.Background()
	var wg sync.WaitGroup

	concurrencyLimit := 5
	sem := make(chan struct{}, concurrencyLimit)

	for _, link := range links {
		wg.Add(1)
		go func(l string) {
			defer wg.Done()
			if err := GetAiData(ctx, s, l, sem); err != nil {
				fmt.Println("Error processing", l, ":", err)
			}
		}(link)
	}
	wg.Wait()
}

func GetAiData(ctx context.Context, s *config.Config, link string, sem chan struct{}) error {
	sem <- struct{}{}        // acquire semaphore slot
	defer func() { <-sem }() // release slot

	// Access your API key as an environment variable
	session, err := NewSession(ctx, s)
	if err != nil {
		return fmt.Errorf("error creating AI client: %v", err)
	}
	fmt.Println("Getting PDF...")
	pdfBytes, err := downloadPDF(link)
	if err != nil {
		return err
	}
	fmt.Println("PDF downloaded.")
	//schema for the AI model
	schema := protocolDataSchema()

	//config
	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema:   schema,
	}

	parts := []*genai.Part{
		genai.NewPartFromText(ai_prompt),
		genai.NewPartFromBytes(pdfBytes, "application/pdf"),
	}

	contents := []*genai.Content{
		genai.NewContentFromParts(parts, genai.RoleUser),
	}

	payload, err := retry(3, 2*time.Second, func() (ProtocolPayload, error) {
		return handleRequest(ctx, *session, session.model, contents, config)
	})

	if err != nil {
		return fmt.Errorf("analyze failed for %s: %w", link, err)
	}

	api.PrintStruct(payload)

	protocol, err := s.Db.CreateProtocolbyScraping(ctx, database.CreateProtocolbyScrapingParams{
		TumorGroup:  payload.ProtocolSummary.TumorGroup,
		Code:        payload.ProtocolSummary.Code,
		Name:        payload.ProtocolSummary.Name,
		Tags:        payload.ProtocolSummary.Tags,
		Notes:       payload.ProtocolSummary.Name,
		RevisedOn:   payload.ProtocolSummary.RevisedOn,
		ActivatedOn: payload.ProtocolSummary.ActivatedOn,
	})
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			// Duplicate key value violation
			protocol, err = s.Db.GetProtocolByCode(ctx, payload.ProtocolSummary.Code)
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
		articleRef, err := s.Db.CreateArticleReference(ctx, database.CreateArticleReferenceParams{
			Title:   article.Title,
			Authors: article.Authors,
			Journal: article.Journal,
			Year:    article.Year,
			Pmid:    article.Pmid,
			Doi:     article.Doi,
		})
		if err != nil {
			if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
				// Duplicate key value violation
				articleRef, err = s.Db.GetArticleReferenceByData(ctx, database.GetArticleReferenceByDataParams{
					Title:   article.Title,
					Authors: article.Authors,
					Journal: article.Journal,
					Year:    article.Year,
				})
				if err != nil {
					fmt.Println("Error getting article ref: ", err)
					return err
				}
			}
		}

		_ = s.Db.AddArticleReferenceToProtocol(ctx, database.AddArticleReferenceToProtocolParams{
			ProtocolID:  protocol.ID,
			ReferenceID: articleRef.ID,
		})
	}

	// Create physicians

	for _, physician := range payload.Physicians {
		email := strings.ToLower(physician.FirstName) + "." + strings.ToLower(physician.LastName) + "@bccancer.bc.ca"
		processedEmail := strings.ReplaceAll(email, " ", "")
		phys, err := s.Db.CreatePhysician(ctx, database.CreatePhysicianParams{
			FirstName: physician.FirstName,
			LastName:  physician.LastName,
			Email:     processedEmail,
			Site:      "vancouver",
		})
		if err != nil {
			if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
				// Duplicate key value violation
				phys, err = s.Db.GetPhysicianByName(ctx, database.GetPhysicianByNameParams{
					FirstName: physician.FirstName,
					LastName:  physician.LastName,
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
			ProtocolID:  protocol.ID,
		})

	}

	// Create protocol eligibility criteria

	for _, eligibility := range payload.ProtocolEligibilityCriteria {
		elig, err := s.Db.InsertEligibilityCriteria(ctx, database.InsertEligibilityCriteriaParams{
			Type:        eligibility.Type,
			Description: eligibility.Description,
		})
		if err != nil {
			if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
				// Duplicate key value violation
				elig, err = s.Db.GetElibilityCriteriaByDescription(ctx, eligibility.Description)

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

	// // Create protocol precautions

	for _, precaution := range payload.ProtocolPrecautions {
		precaut, err := s.Db.CreateProtocolPrecaution(ctx, database.CreateProtocolPrecautionParams{
			Title:       precaution.Title,
			Description: precaution.Description,
		})
		if err != nil {
			if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
				// Duplicate key value violation
				precaut, err = s.Db.GetProtocolPrecautionByTitleAndDescription(ctx, database.GetProtocolPrecautionByTitleAndDescriptionParams{
					Title:       precaution.Title,
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
			ProtocolID:   protocol.ID,
			PrecautionID: precaut.ID,
		})
	}

	// // Create Protocol Cautions

	for _, caution := range payload.ProtocolCautions {
		caut, err := s.Db.CreateProtocolCaution(ctx, caution.Description)
		if err != nil {
			if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
				// Duplicate key value violation
				caut, err = s.Db.GetProtocolCautionByDescription(ctx, caution.Description)

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
			CautionID:  caut.ID,
		})
	}

	// // Create Tests

	for _, testGroup := range payload.TestGroups {

		test_category, err := s.Db.UpsertProtoTestCategory(ctx, database.UpsertProtoTestCategoryParams{
			ID:         testGroup.ID,
			ProtocolID: protocol.ID,
			Category:   testGroup.Category,
			Comments:   testGroup.Comments,
			Position:   testGroup.Position,
		})
		if err != nil {
			if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
				// Duplicate key value violation
				response, err := s.Db.GetTestCategoryByID(ctx, testGroup.ID)
				if err != nil {
					fmt.Println("Error getting test category: ", err)
					return err
				}
				test_category, err = api.ToResponseData[database.ProtocolTest](response)
				if err != nil {
					fmt.Println("Error converting test category: ", err)
					return err
				}

			} else {
				fmt.Println("Error creating test category: ", err)
				return err
			}
		}

		for _, test := range testGroup.Tests {
			added, err := s.Db.UpsertTest(ctx, database.UpsertTestParams{
				ID:           api.ParseOrNilUUID(""),
				Name:         test.Name,
				Description:  test.Description,
				FormUrl:      "",
				Unit:         "",
				LowerLimit:   0,
				UpperLimit:   0,
				TestCategory: "",
			})
			if err != nil {

				if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
					// Duplicate key value violation
					added, err = s.Db.GetTestByName(ctx, test.Name)
					if err != nil {
						fmt.Println("Error getting test: ", err)
						return err
					}
				} else {
					fmt.Println("Error creating test: ", err)
					return err
				}
			}

			err = s.Db.AddTestToProtoTestCategory(ctx, database.AddTestToProtoTestCategoryParams{
				TestsID:         added.ID,
				ProtocolTestsID: test_category.ID,
			})
			if err != nil {
				return err
			}
		}
	}

	// // Create Prescriptions

	for _, medGroup := range payload.PrescriptionGroups {

		test_category, err := s.Db.UpsertProtoMedCategory(ctx, database.UpsertProtoMedCategoryParams{
			ID:         medGroup.ID,
			ProtocolID: protocol.ID,
			Category:   medGroup.Category,
			Comments:   medGroup.Comments,
		})
		if err != nil {
			if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
				// Duplicate key value violation
				response, err := s.Db.GetMedCategoryByID(ctx, medGroup.ID)
				if err != nil {
					fmt.Println("Error getting med category: ", err)
					return err
				}
				test_category, err = api.ToResponseData[database.ProtocolMed](response)
				if err != nil {
					fmt.Println("Error converting med category: ", err)
					return err
				}

			} else {
				fmt.Println("Error creating med category: ", err)
				return err
			}
		}

		for _, px := range medGroup.Prescriptions {

			med, err := s.Db.AddMedication(ctx, database.AddMedicationParams{
				Name:           px.MedicationName,
				Description:    px.MedicationDescription,
				Category:       px.MedicationCategory,
				AlternateNames: px.MedicationAlternates,
			})
			if err != nil {
				if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
					// Duplicate key value violation
					med, err = s.Db.GetMedicationByName(ctx, px.MedicationName)
					if err != nil {
						fmt.Println("Error getting medication: ", err)
						return err
					}
				}
			}

			added, err := s.Db.UpsertPrescription(ctx, database.UpsertPrescriptionParams{
				ID:           api.ParseOrNilUUID(""),
				MedicationID: med.ID,
				Dose:         px.Dose,
				Route:        px.Route,
				Frequency:    px.Frequency,
				Duration:     px.Duration,
				Instructions: px.Instructions,
				Renewals:     px.Renewals,
			})
			if err != nil {

				if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
					// Duplicate key value violation
					added.ID, err = s.Db.GetPrescriptionsByArguments(ctx, database.GetPrescriptionsByArgumentsParams{
						MedicationID: med.ID,
						Dose:         px.Dose,
						Route:        px.Route,
						Frequency:    px.Frequency,
						Duration:     px.Duration,
						Instructions: px.Instructions,
					})
					if err != nil {
						fmt.Println("Error getting prescription: ", err)
						return err
					}
				} else {
					fmt.Println("Error creating prescription: ", err)
					return err
				}
			}

			err = s.Db.AddPrescriptionToProtocolCategory(ctx, database.AddPrescriptionToProtocolCategoryParams{
				MedicationPrescriptionID: added.ID,
				ProtocolMedsID:           test_category.ID,
			})

			if err != nil {
				return err
			}
		}
	}

	// // Create cycles

	for _, cycle := range payload.ProtocolCycles {

		added_cycle, err := s.Db.UpsertCycleToProtocol(ctx, database.UpsertCycleToProtocolParams{
			ID:            cycle.ID,
			ProtocolID:    protocol.ID,
			Cycle:         cycle.Cycle,
			CycleDuration: cycle.CycleDuration,
		})
		if err != nil {
			if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
				// Duplicate key value violation
				added_cycle, err = s.Db.GetCycleByData(ctx, database.GetCycleByDataParams{
					Cycle:         cycle.Cycle,
					CycleDuration: cycle.CycleDuration,
					ProtocolID:    protocol.ID,
				})
				if err != nil {
					fmt.Println("Error getting protocol cycle: ", err)
					return err
				}

			} else {
				fmt.Println("Error creating protocol cycle: ", err)
				return err
			}
		}

		for _, tx := range cycle.Treatments {

			med, err := s.Db.AddMedication(ctx, database.AddMedicationParams{
				Name:           tx.MedicationName,
				Description:    tx.MedicationDescription,
				Category:       tx.MedicationCategory,
				AlternateNames: tx.MedicationAlternates,
			})
			if err != nil {
				if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
					// Duplicate key value violation
					med, err = s.Db.GetMedicationByName(ctx, tx.MedicationName)
					if err != nil {
						fmt.Println("Error getting medication: ", err)
						return err
					}
				}
			}

			added, err := s.Db.UpsertProtocolTreatment(ctx, database.UpsertProtocolTreatmentParams{
				ID:                  api.ParseOrNilUUID(""),
				MedicationID:        med.ID,
				Dose:                tx.Dose,
				Route:               tx.Route,
				Frequency:           tx.Frequency,
				Duration:            tx.Duration,
				AdministrationGuide: tx.AdministrationGuide,
			})

			if err != nil {

				if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
					// Duplicate key value violation
					added, err = s.Db.GetProtocolTreatmentByData(ctx, database.GetProtocolTreatmentByDataParams{
						MedicationID: med.ID,
						Dose:         tx.Dose,
						Route:        tx.Route,
						Frequency:    tx.Frequency,
						Duration:     tx.Duration,
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

			err = s.Db.AddTreatmentToCycle(ctx, database.AddTreatmentToCycleParams{
				ProtocolCyclesID:    added_cycle.ID,
				ProtocolTreatmentID: added.ID,
			})
			if err != nil {
				return err
			}
		}
	}

	// // Create Toxicity Modifications

	for _, tox := range payload.Toxicities {

		toxicity, err := s.Db.AddToxicity(ctx, database.AddToxicityParams{
			Title:       tox.Title,
			Description: tox.Description,
			Category:    tox.Category,
		})

		if err != nil {
			if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
				// Duplicate key value violation
				toxicity, err = s.Db.GetToxicityByName(ctx, tox.Title)

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
			grade, err := s.Db.AddToxicityGrade(ctx, database.AddToxicityGradeParams{
				Grade:       database.GradeEnum(mod.Grade),
				Description: mod.GradeDescription,
				ToxicityID:  toxicity.ID,
			})
			if err != nil {
				if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
					// Duplicate key value violation
					grade, err = s.Db.GetToxicityGradeByGrade(ctx, database.GetToxicityGradeByGradeParams{
						Grade:      database.GradeEnum(mod.Grade),
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

			_, err = s.Db.AddToxicityModification(ctx, database.AddToxicityModificationParams{
				Adjustment:      mod.Adjustment,
				ToxicityGradeID: grade.ID,
				ProtocolID:      protocol.ID,
			})

			if err != nil {
				fmt.Println("Error creating toxicity modification: ", err)
			}
		}

	}	

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
func handleRequest(ctx context.Context, s Session, model string, contents []*genai.Content, config *genai.GenerateContentConfig) (payload ProtocolPayload, err error) {

	response, err := s.client.Models.GenerateContent(ctx, model, contents, config)
	if err != nil {
		fmt.Printf("Request failed: %v\n", err)
		return
	}

	for _, c := range response.Candidates {
		for _, part := range c.Content.Parts {
			// Parse the extracted JSON data
			err = json.Unmarshal([]byte(part.Text), &payload)
			if err != nil {
				fmt.Printf("Request - Error unmarshaling JSON: %v\n", err)
				return ProtocolPayload{}, err
			}
		}
	}

	return payload, nil
}
