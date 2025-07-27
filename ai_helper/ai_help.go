package ai_helper

import (
	// "bcca_crawler/api"
	"bcca_crawler/internal/config"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"google.golang.org/genai"
)

// Function to handle each request, send it, and parse the response
func handle_Request(ctx context.Context, s Session1, model string, contents []*genai.Content, config *genai.GenerateContentConfig) (payload ProtocolPayload, err error) {

	response, err := s.client.Models.GenerateContent(ctx, model, contents, config)
	if err != nil {
		fmt.Printf("Request failed: %v\n", err)
		return
	}

	fmt.Printf("Response received:%v\n", response.UsageMetadata.TotalTokenCount)

	for _, c := range response.Candidates {
		for _, part := range c.Content.Parts {
				fmt.Println("Raw AI FileData response:", part.FileData )
				fmt.Println("Raw AI InlineData response:", part.InlineData )			
				fmt.Println("Raw AI text response:", part.Text)
			
		}
	}

	jsonbytes, err := response.MarshalJSON()
	if err != nil {
		fmt.Printf("Error marshalling JSON failed: %v\n", err)
		return ProtocolPayload{}, err
	}

	// Parse the extracted JSON data
	err = json.Unmarshal([]byte(jsonbytes), &payload)
	if err != nil {
		fmt.Printf("Request - Error unmarshaling JSON: %v\n", err)
		return ProtocolPayload{}, err
	}

	return payload, nil
}

// Function to download PDF content
func download_PDF(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

const ai_prompt1 = `You are a medical AI model tasked with analyzing the joined PDF document and extracting structured information in JSON format.
### Task:
1. Parse the provided PDF thoroughly.
2. Extract all relevant information
3. toxicities should be extracted using the CTCAE v5 terminology.
4. All the tests should be single word names, and if there are multiple tests, they should be a new entry in the array.
5. Return the completed JSON object, ensuring all fields are validated for data type and consistency.`

type Session1 struct {
	ctx    context.Context
	client *genai.Client
	model  string
}

func New_Session(ctx context.Context, s *config.Config) (*Session1, error) {
	client, err := genai.NewClient(ctx,
		&genai.ClientConfig{
			APIKey:  s.GeminiApiKey,
			Backend: genai.BackendGeminiAPI,
		})
	if err != nil {
		return nil, fmt.Errorf("error creating AI client: %v", err)
	}

	model := "gemini-2.5-flash"

	return &Session1{
		ctx:    ctx,
		client: client,
		model:  model,
	}, nil
}

func Get_AiData(ctx context.Context, s *config.Config, link string, sem chan struct{}) error {
	sem <- struct{}{}        // acquire semaphore slot
	defer func() { <-sem }() // release slot

	// Access your API key as an environment variable
	session, err := New_Session(ctx, s)
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
		// ResponseMIMEType: "application/json",
		ResponseSchema:   schema,
	}

	parts := []*genai.Part{
		genai.NewPartFromText(ai_prompt),
		genai.NewPartFromBytes(pdfBytes, "application/pdf"),		
	}

	contents := []*genai.Content{
		genai.NewContentFromParts(parts, genai.RoleUser),
	}

	payload, err := handle_Request(ctx, *session, session.model, contents, config)
	// Retry block for the AI payload
	// payload, err := retry(3, 2*time.Second, func() (ProtocolPayload, error) {
	// 	return handleRequest(ctx, model, req)
	// })

	if err != nil {
		return fmt.Errorf("analyze failed for %s: %w", link, err)
	}

	fmt.Println("Payload:")
	// api.PrintStruct(payload)
	fmt.Printf("%+v\n", payload)
	return nil
}
