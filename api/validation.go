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

// Custom validation function
func TumorGroupValidator(fl validator.FieldLevel) bool {
	tumorGroup := strings.ToLower(fl.Field().String()) // Ensure case-insensitivity
	return validTumorGroups[tumorGroup]
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
