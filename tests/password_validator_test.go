package main

import (
	"testing"
	"github.com/go-playground/validator/v10"
	"bcca_crawler/api"	
)

func TestPasswordStrengthValidator(t *testing.T) {
	// Create a mock field level object for testing
	validate := validator.New()
	validate.RegisterValidation("passwordstrength", api.PasswordStrengthValidator)



	tests := []struct {
		password    api.CreateUserRequest
		expectedValid bool
		expectedErrors []string
	}{
		{
			password:      api.CreateUserRequest{Email: "test@email.com", Password: "StrongP@ssw0rd"},
			expectedValid: true,
			expectedErrors: nil,
		},
		{
			password:      api.CreateUserRequest{Email: "test@email.com", Password: "short"}, 
			expectedValid: false,
			expectedErrors: []string{
				"Password must be at least 8 characters long.",
				"Password must contain at least one uppercase letter.",
				"Password must contain at least one number.",
				"Password must contain at least one special character.",
			},
		},
		{
			password:      api.CreateUserRequest{Email: "test@email.com", Password: "short"}, 
			expectedValid: true,
			expectedErrors: nil,
		},
		{
			password:      api.CreateUserRequest{Email: "test@email.com", Password: "password123"}, 
			expectedValid: false,
			expectedErrors: []string{
				"Password must contain at least one uppercase letter.",
				"Password must contain at least one special character.",
			},
		},
		{
			password:      api.CreateUserRequest{Email: "test@email.com", Password: "PASSWORD123"}, 
			expectedValid: false,
			expectedErrors: []string{
				"Password must contain at least one lowercase letter.",
				"Password must contain at least one special character.",
			},
		},
		{
			password:      api.CreateUserRequest{Email: "test@email.com", Password: "Password@"}, 
			expectedValid: false,
			expectedErrors: []string{
				"Password must contain at least one number.",
			},
		},
		{
			password:      api.CreateUserRequest{Email: "test@email.com", Password: "Pass123"}, 
			expectedValid: false,
			expectedErrors: []string{
				"Password must be at least 8 characters long.",
				"Password must contain at least one special character.",
			},
		},
	}

	// Loop through the test cases
	for _, tt := range tests {
		t.Run(tt.password.Password, func(t *testing.T) {			

			// Call PasswordStrengthValidator
			err := validate.Struct(tt.password)

			// Check if the result matches expectations
			if err != nil && tt.expectedValid {
				t.Errorf("expected valid: %v, got: %v", tt.expectedValid, false)				
			}

			// Check if the error messages match expectations
			if (err == nil && tt.expectedValid == false)  {
				t.Errorf("expected valid: %v, got: %v", tt.expectedValid, true)
			}

		})
	}
}
