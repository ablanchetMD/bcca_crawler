package api

import (
	"bcca_crawler/internal/config"
	"github.com/go-playground/validator/v10"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"github.com/google/uuid"
	"database/sql"
	"reflect"
	"strings"
	
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

type IDs struct {
	ID uuid.UUID
	ProtocolID uuid.UUID
	CycleID uuid.UUID
}

func ParseAndValidateID(r *http.Request) (IDs, error) {
	var ids IDs    
	
	if id := r.PathValue("id"); id != "" {
		parsed_id, err := uuid.Parse(id)
		if err != nil {
			return ids, fmt.Errorf("id is not a valid uuid")
		}
		ids.ID = parsed_id
	}

	if protocol_id := r.PathValue("protocol_id"); protocol_id != "" {
        pid, err := uuid.Parse(protocol_id)
        if err != nil {
            return ids, fmt.Errorf("protocol_id is not a valid uuid")
        }
        ids.ProtocolID = pid
    }

	if cycle_id := r.PathValue("cycle_id"); cycle_id != "" {
		cid, err := uuid.Parse(cycle_id)
		if err != nil {
			return ids, fmt.Errorf("cycle_id is not a valid uuid")
		}
		ids.CycleID = cid
	}
    
    return ids, nil
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
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors := make(map[string]string)
			for _, e := range validationErrors {
				errors[e.Field()] = fmt.Sprintf("failed validation: %s", e.Tag())
			}
			return fmt.Errorf("validation failed: %v", errors)
		}
		return fmt.Errorf("validation failed: %s", err.Error())
	}
	return nil
}

func AnalyzeTreeWithCustomTypes(
	data interface{},
	typeMap map[string][]interface{},
	customTypes map[string]reflect.Type,
) {
	switch node := data.(type) {
	case []interface{}: // Handle arrays
		for _, item := range node {
			AnalyzeTreeWithCustomTypes(item, typeMap, customTypes)
		}
	case map[string]interface{}: // Handle objects (nodes)
		for _, v := range node {
			AnalyzeTreeWithCustomTypes(v, typeMap, customTypes)
		}

		// Try to match this object to a custom type
		for typeName, typeDef := range customTypes {
			matchedObj := reflect.New(typeDef).Interface()
			bytes, err := json.Marshal(node)
			if err != nil {
				continue
			}
			err = json.Unmarshal(bytes, matchedObj)
			if err == nil {
				// Match found
				typeMap[typeName] = append(typeMap[typeName], matchedObj)
			}
		}
	default:
		// Handle primitive types
		typeName := fmt.Sprintf("%T", node)
		typeMap[typeName] = append(typeMap[typeName], node)
	}
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

// Return an array of tuples into an array of structs
func ConvertTuplesToStructs[T any](t interface{},fields ...string) ([]T,error) {
	var returnArray []T

	convertByte,ok := t.([]byte)
	if !ok {
		return nil, fmt.Errorf("invalid protocol data format: expected []byte, got %T", t)
	}		

	tuplesString := string(convertByte)	

	// Remove the surrounding parentheses and split by "),("
	tuplesString = tuplesString[3 : len(tuplesString)-3] // Remove outer parentheses
	
	// Split into individual tuples
	tuples := strings.Split(tuplesString, "),(")	
	for _, tuple := range tuples {		
		// Split each tuple by comma
		parts := strings.Split(tuple, ",")		
		
		var data T
		dataVal := reflect.ValueOf(&data).Elem()

		structType := reflect.TypeOf(data)
		if len(parts) != structType.NumField() {					
				return nil, fmt.Errorf("invalid protocol data format: expected %d fields, got %d", len(fields), structType.NumField())
			}		
		// Loop through the fields of the struct and assign values
		for i := 0; i < structType.NumField(); i++ {
			// Get the struct field by index
			field := structType.Field(i)

			// Check if the field has a tag and whether it's valid for assignment
			if field.Tag.Get("json") != "" {
				// Get the field value by index
				fieldVal := dataVal.Field(i)
				if fieldVal.IsValid() && fieldVal.CanSet() {
					// Set the value of the field
					fieldVal.SetString(strings.TrimSpace(parts[i]))
				}
			}
		}
		returnArray = append(returnArray, data)			
		}

	return returnArray,nil
}

	


func PrintStruct(v interface{}) {
	// Convert the struct to pretty JSON
	jsonData, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling struct: %v", err)
		return
	}

	// Print the JSON to the console
	fmt.Println(string(jsonData))
}

