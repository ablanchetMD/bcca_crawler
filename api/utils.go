package api

import (
	"bcca_crawler/internal/config"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"github.com/google/uuid"
	"database/sql"
	
)

type QueryParams struct {
	Sort    string
	SortBy  string
	Page    int
	Limit   int
	Offset  int
	FilterBy  string
	Fields  []string
	Include []string
	Exclude []string
}

type IPApiResponse struct {
	Query string `json:"query"`
	Country string `json:"country"`
	City string `json:"city"`
	Region string `json:"region"`
	Zip string `json:"zip"`
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
	Timezone string `json:"timezone"`
	ISP string `json:"isp"`
	Org string `json:"org"`
	AS string `json:"as"`
}

func ParseAndValidateID(r *http.Request) (uuid.UUID, error) {
    id := r.PathValue("id")
    if len(id) == 0 {
        return uuid.Nil, fmt.Errorf("no id provided")
    }

    parsed_id, err := uuid.Parse(id)
    if err != nil {
        return uuid.Nil, fmt.Errorf("id is not a valid uuid")
    }
    return parsed_id, nil
}

func UnmarshalAndValidatePayload(c *config.Config,r *http.Request, v interface{}) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("no body in request")
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &v)
	if err != nil {
		return fmt.Errorf("invalid request body")
	}	
	
	// Validate the request data
	err = c.Validate.Struct(v)
	if err != nil {		
		return fmt.Errorf("validation failed: %s", err.Error())
	}
	return nil
}

func GetIPGeoLocation(ip string) (IPApiResponse, error) {	
	var ipApiResponse IPApiResponse
	// Make a request to ip-api.com
	rawURL := fmt.Sprintf("http://ip-api.com/json/%s", ip)

	resp, err := http.Get(rawURL)

	if err != nil {
		return ipApiResponse, err
	}
	if resp.StatusCode != http.StatusOK {
		return ipApiResponse, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ipApiResponse, fmt.Errorf("no body in request")
	}	

	defer resp.Body.Close()	

	err = json.Unmarshal(body, &ipApiResponse)
	if err != nil {
		return ipApiResponse, fmt.Errorf("invalid request body")
	}

	return ipApiResponse, nil
}

// Convert a string to sql.NullString
func ToNullString(s string) sql.NullString {
    return sql.NullString{
        String: s,
        Valid:  s != "", // Set Valid to true if the string is non-empty
    }
}

