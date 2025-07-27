package api

import (
	"bcca_crawler/internal/config"
	"bcca_crawler/internal/json_utils"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type QueryParams struct {
	Sort     string
	SortBy   string
	Page     int
	Limit    int
	Offset   int
	FilterBy string
	Fields   []string
	Include  []string
	Exclude  []string
}

type HandlerError struct {
	StatusCode int
	Message    string        // for user-facing response
	Err        error         // for internal logs
}

func (e *HandlerError) Error() string {
	return fmt.Sprintf("status: %d, message: %s, internal: %v", e.StatusCode, e.Message, e.Err)
}

func WrapError(code int, msg string, err error) *HandlerError {
	return &HandlerError{StatusCode: code, Message: msg, Err: err}
}

type IPApiResponse struct {
	Query    string  `json:"query"`
	Country  string  `json:"country"`
	City     string  `json:"city"`
	Region   string  `json:"region"`
	Zip      string  `json:"zip"`
	Lat      float64 `json:"lat"`
	Lon      float64 `json:"lon"`
	Timezone string  `json:"timezone"`
	ISP      string  `json:"isp"`
	Org      string  `json:"org"`
	AS       string  `json:"as"`
}

type IDs struct {
	ID         uuid.UUID
	ProtocolID uuid.UUID
	CycleID    uuid.UUID
	PxCategoryID uuid.UUID
	LabCategoryID uuid.UUID	
}

type MapperFunc[T any] func(dbResult any) (T, error)

type UpsertFunc[Req any] func(cfg *config.Config, ctx context.Context, req Req, ids IDs) error
type GetFunc[Res any] func(cfg *config.Config, ctx context.Context, ids IDs) (Res, error)
type GetFuncWithQueries[Res any] func(cfg *config.Config, ctx context.Context, ids IDs, query url.Values) (Res, error)
type ModifierFunc func(cfg *config.Config, ctx context.Context, ids IDs) (string, error)

func ParseAndValidateID(r *http.Request) (IDs, error) {
	var ids IDs
	vars := mux.Vars(r)
	if id := vars["id"]; id != "" {
		parsed_id, err := uuid.Parse(id)
		if err != nil {
			return ids, fmt.Errorf("id is not a valid uuid")
		}
		ids.ID = parsed_id
	}

	if protocol_id := vars["protocol_id"]; protocol_id != "" {
		pid, err := uuid.Parse(protocol_id)
		if err != nil {
			return ids, fmt.Errorf("protocol_id is not a valid uuid")
		}
		ids.ProtocolID = pid
	}

	if cycle_id := vars["cycle_id"]; cycle_id != "" {
		cid, err := uuid.Parse(cycle_id)
		if err != nil {
			return ids, fmt.Errorf("cycle_id is not a valid uuid")
		}
		ids.CycleID = cid
	}

	if px_category_id := vars["px_category_id"]; px_category_id != "" {
		pxid, err := uuid.Parse(px_category_id)
		if err != nil {
			return ids, fmt.Errorf("px_category_id is not a valid uuid")
		}
		ids.PxCategoryID = pxid
	}

	if lab_category_id := vars["lab_category_id"]; lab_category_id != "" {
		lcid, err := uuid.Parse(lab_category_id)
		if err != nil {
			return ids, fmt.Errorf("lab_category_id is not a valid uuid")
		}
		ids.LabCategoryID = lcid
	}

	return ids, nil
}

func ParseOrGenerateUUID(s string) uuid.UUID {
	if u, err := uuid.Parse(s); err == nil {
		return u
	}
	return uuid.New()
}

func ParseOrNilUUID(s string) uuid.UUID {
	if u, err := uuid.Parse(s); err == nil {
		return u
	}
	return uuid.Nil
}

func HandleUpsert[Req any](
	c *config.Config,
	w http.ResponseWriter,
	r *http.Request,
	upsertFn UpsertFunc[Req],
) {
	var req Req
	fmt.Println("Upsert")
	ids, err := ParseAndValidateID(r)
	if err != nil {		
		fmt.Println("Parsing error")
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := UnmarshalAndValidatePayload(c, r, &req); err != nil {
		fmt.Printf("Unmarshall error:%s",err)
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = upsertFn(c, r.Context(), req, ids)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	json_utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Item upserted successfully"})
}

func HandleGet[Res any](
	c *config.Config,
	w http.ResponseWriter,
	r *http.Request,
	getFn GetFunc[Res],
) {
	ids, err := ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := getFn(c, r.Context(), ids)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	json_utils.RespondWithJSON(w, http.StatusOK, res)
}

func HandleGetWithQ[Res any](
	c *config.Config,
	w http.ResponseWriter,
	r *http.Request,
	getFn GetFuncWithQueries[Res],
) {
	ids, err := ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	query := r.URL.Query()

	res, err := getFn(c, r.Context(), ids,query)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	json_utils.RespondWithJSON(w, http.StatusOK, res)
}

func HandleModify(
	c *config.Config,
	w http.ResponseWriter,
	r *http.Request,
	modifierFn ModifierFunc,
) {	
	ids, err := ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	msg, err := modifierFn(c, r.Context(), ids)
	if err != nil {
		fmt.Printf("Handle Modify Error with Request: %v\n", r.URL)
		json_utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	json_utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": msg})
}

func UnmarshalAndValidatePayload(c *config.Config, r *http.Request, v interface{}) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("no body in request")
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &v)
	if err != nil {
		fmt.Printf("request body : %s,\n error: %s",body,err)
		return fmt.Errorf("invalid request body")
	}

	// Validate the request data
	err = c.Validate.Struct(v)
	if err != nil {
		PrintStruct(v)
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
func ConvertTuplesToStructs[T any](t interface{}, fields ...string) ([]T, error) {
	var returnArray []T

	convertByte, ok := t.([]byte)
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

	return returnArray, nil
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
