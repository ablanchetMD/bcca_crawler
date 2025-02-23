package api

import(
	"github.com/go-playground/validator/v10"
	"unicode"

	"strings"
)

// Predefined list of valid tumor group codes
var validTumorGroups = map[string]bool{
	"lymphoma": true,
	"myeloma": true,
	"bmt": true,
	"leukemia": true,
	"breast": true,
	"gastrointestinal": true,
	"genitourinary": true,
	"gynecology": true,
	"head&neck": true,
	"lung": true,
	"melanoma": true,
	"neuro-oncology": true,
	"sarcoma": true,	
	// Add more as needed
}

var validEligiblityCriteria = map[string]bool{
	"inclusion": true,
	"exclusion": true,
	"notes": true,
	"unknown": true,
}

var validTestProtocolCategories = map[string]bool{
	"baseline": true,
	"followup": true,
	"unknown": true,
}

var validTestProtocolUrgency = map[string]bool{
	"urgent": true,
	"non_urgent": true,
	"if_necessary": true,
	"unknown": true,
}

var validProtocolPrescriptionCategory = map[string]bool{
	"premed": true,
	"support": true,
	"unknown": true,
}

var validPrescriptionRoutes = map[string]bool{
	"oral": true,
	"iv": true,
	"im": true,
	"sc": true,
	"topical": true,
	"unknown": true,
}

var validGrades = map[string]bool{
	"1": true,
	"2": true,
	"3": true,
	"4": true,
	"unknown": true,
}

var validPhysicianSites = map[string]bool{
	"vancouver": true,
	"victoria": true,
	"abbotsford": true,
	"kelowna": true,
	"prince_george": true,
	"nanaimo": true,
	"surrey": true,
	"unknown": true,
}


// Custom validation function
func TumorGroupValidator(fl validator.FieldLevel) bool {
	tumorGroup := strings.ToLower(fl.Field().String()) // Ensure case-insensitivity
	return validTumorGroups[tumorGroup]
}

func GradeValidator(fl validator.FieldLevel) bool {
	grade := strings.ToLower(fl.Field().String()) // Ensure case-insensitivity
	return validGrades[grade]
}

func ProtocolPrescriptionCategoryValidator(fl validator.FieldLevel) bool {
	prescription_category := strings.ToLower(fl.Field().String()) // Ensure case-insensitivity
	return validProtocolPrescriptionCategory[prescription_category]
}

func EligibilityCriteriaValidator(fl validator.FieldLevel) bool {
	eligibilityCriteria := strings.ToLower(fl.Field().String()) // Ensure case-insensitivity
	return validEligiblityCriteria[eligibilityCriteria]
}

func PrescriptionRouteValidator(fl validator.FieldLevel) bool {
	prescription_route := strings.ToLower(fl.Field().String()) // Ensure case-insensitivity
	return validPrescriptionRoutes[prescription_route]
}

func TestCategoryValidator(fl validator.FieldLevel) bool {
	testCategory := strings.ToLower(fl.Field().String()) // Ensure case-insensitivity
	return validTestProtocolCategories[testCategory]
}

func TestUrgencyValidator(fl validator.FieldLevel) bool {
	testUrgency := strings.ToLower(fl.Field().String()) // Ensure case-insensitivity
	return validTestProtocolUrgency[testUrgency]
}

func PhysicianSiteValidator(fl validator.FieldLevel) bool {
	physicianSite := strings.ToLower(fl.Field().String()) // Ensure case-insensitivity
	return validPhysicianSites[physicianSite]
}

// PasswordStrengthValidator checks for strong passwords using bitwise operations.
func PasswordStrengthValidator(fl validator.FieldLevel) (bool) {
	password := fl.Field().String()
	
	const (
		hasLower  = 1 << iota // 0001
		hasUpper              // 0010
		hasNumber             // 0100
		hasSpecial            // 1000
	)
	requiredMask := hasLower | hasUpper | hasNumber | hasSpecial // 1111

	var flags int
	length := len(password)

	if length < 8 { // Quickly fail for short passwords
		return false
	}

	if length > 50 { // Quickly fail for long passwords
		return false
	}

	for _, char := range password {
		switch {
		case unicode.IsLower(char):
			flags |= hasLower
		case unicode.IsUpper(char):
			flags |= hasUpper
		case unicode.IsDigit(char):
			flags |= hasNumber
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			flags |= hasSpecial
		}

		// Early exit if all criteria are met
		if flags == requiredMask {
			return true
		}
	}	

	// Check if all criteria are met
	return flags == requiredMask
}
