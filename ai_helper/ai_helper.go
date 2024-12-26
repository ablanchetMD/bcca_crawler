package ai_helper

import (
	"context"
	"fmt"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
	"io"
	"net/http"

	"bcca_crawler/internal/config"

)

type ProtocolSummary struct {
	Title            string   `json:"title"`
	ProtocolCode     string   `json:"protocol_code"`
	TumourGroup      string   `json:"tumour_group"`
	ContactPhysicians []string `json:"contact_physicians"`
}

type Eligibility struct {
	InclusionCriteria []string `json:"inclusion_criteria"`
	Notes    []string `json:"notes"`
	ExclusionCriteria []string `json:"exclusion_criteria"`
}

type Tests struct {
	Baseline struct {
		RequiredBeforeTreatment []string `json:"required_before_treatment"`
		RequiredButCanProceed   []string `json:"required_but_can_proceed"`
		IfClinicallyIndicated   []string `json:"if_clinically_indicated"`
	} `json:"baseline"`
	FollowUp            []string `json:"follow_up"`
	IfClinicallyIndicated []string `json:"if_clinically_indicated"`
}

type ToxicityManagement struct {
	Grade string 	`json:"grade"`
	Management string `json:"adjustment"`
}

type Toxicity struct {
	Toxicity   string `json:"toxicity"`
	ToxicityManagement []ToxicityManagement `json:"toxicity_management"`
}

type DoseModifications struct {
	Toxicity           []Toxicity  `json:"dose_modifications"`
	HepaticImpairment struct {
		MildModerate string `json:"mild_moderate"`
		Severe       string `json:"severe"`
	} `json:"hepatic_impairment"`
	RenalImpairment struct {
		MildModerate string `json:"mild_moderate"`
		Severe       string `json:"severe"`
	} `json:"renal_impairment"`
}

type Treatment struct {
	Drug           string `json:"drug"`
	Dose           string `json:"dose"`
	Administration string `json:"administration"`
	Duration       string `json:"duration"`
}

type Reference struct {
	Title       string `json:"title"`
	Authors     string `json:"authors"`
	Journal     string `json:"journal,omitempty"`
	Year        int    `json:"year,omitempty"`	
}

type ProtocolData struct {
	ProtocolSummary    ProtocolSummary     `json:"protocol_summary"`
	Eligibility        Eligibility         `json:"eligibility"`	
	Cautions           []string            `json:"cautions"`
	Tests              Tests               `json:"tests"`
	Treatment          Treatment           `json:"treatment"`
	DoseModifications  DoseModifications   `json:"dose_modifications"`
	Precautions        []string            `json:"precautions"`	
	References         []Reference         `json:"references"`
}

const ai_prompt = `You are an AI tasked with extracting structured data from PDFs. Extract the relevant information from the PDF and return it as JSON matching the structure below: 

{
  "protocol_summary": {
    "title": "string",
    "protocol_code": "string",
    "tumour_group": "string",
    "contact_physicians": ["string"]
  },
  "eligibility": {
    "inclusion_criteria": ["string"],
    "notes": ["string"],
    "exclusion_criteria": ["string"]
  },
  "tests": {
    "baseline": {
      "required_before_treatment": ["string"],
      "required_but_can_proceed": ["string"],
      "if_clinically_indicated": ["string"]
    },
    "follow_up": ["string"],
    "if_clinically_indicated": ["string"]
  },
  "treatment": {
    "drug": "string",
    "dose": "string",
    "administration": "string",
    "duration": "string"
  },
  "dose_modifications": {
    "dose_modifications": [
      {
        "toxicity": "string",
        "toxicity_management": [
          {
            "grade": "string",
            "adjustment": "string"
          }
        ]
      }
    ],
    "hepatic_impairment": {
      "mild_moderate": "string",
      "severe": "string"
    },
    "renal_impairment": {
      "mild_moderate": "string",
      "severe": "string"
    }
  },
  "cautions": ["string"],
  "precautions": ["string"],
  "references": [
    {
      "title": "string",
      "authors": "string",
      "journal": "string",
      "year": "int"
    }
  ]
}

Instructions: 
1. Analyze the PDF to identify sections matching the fields in the data structure. 
2. Extract and format the data into JSON, ensuring all fields are populated according to the structure, even if some fields are empty. 
3. For lists (e.g., contact_physicians, inclusion_criteria), ensure each item in the list is an individual string. 
4. For subfields like hepatic_impairment, map the information directly from the relevant PDF sections.

If there is ambiguity or missing data, indicate this with "unknown" or an empty string for strings and [] for lists.

Expected Output Example: 
Provide the extracted data as JSON adhering strictly to the structure above.`

  

func TestAi(s *config.Config) error {
	ctx := context.Background()
	// Access your API key as an environment variable
	client, err := genai.NewClient(ctx, option.WithAPIKey(s.GeminiApiKey))
	if err != nil {
		return fmt.Errorf("error creating AI client: %v", err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")

	 // Download the pdf.
	 pdfResp, err := http.Get("http://www.bccancer.bc.ca/chemotherapy-protocols-site/Documents/Lymphoma-Myeloma/LYACAL_Protocol.pdf")
	 if err != nil {
		 return err
	 }
	 defer pdfResp.Body.Close()
 
	 pdfBytes, err := io.ReadAll(pdfResp.Body)
	 if err != nil {
		 return err
	 }
 
	 // Create the request.
	 req := []genai.Part{
		 genai.Blob{MIMEType: "application/pdf", Data: pdfBytes},
		 genai.Text(ai_prompt),
 
		//  genai.Text("Generate a JSON file with this following format:"+jsonFormat+" based on the content of the PDF."),
	 }
 
	 // Generate content.
	 resp, err := model.GenerateContent(ctx, req...)
	 if err != nil {
		 return err
	 }
 
	 // Handle the response of generated text.
	 for _, c := range resp.Candidates {
		 if c.Content != nil {
			 fmt.Println(*c.Content)
		 }
	 }
	return nil
}
